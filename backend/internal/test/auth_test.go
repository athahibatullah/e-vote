package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"e-vote-system/config"
	"e-vote-system/internal/auth"
	"e-vote-system/internal/handler"
	"e-vote-system/internal/mail"
	"e-vote-system/middleware"
	database "e-vote-system/models/database"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/auth/login", auth.Login)
	r.Post("/api/auth/verify-otp", auth.VerifyOTP)
	return r
}

func prepareTestDB() {
	// ID := uuid.New()
	// Fake user for testing
	pasword, err := bcrypt.GenerateFromPassword([]byte("abcde"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	var testVoter = database.Voter{
		ID:         "",
		VoterName:  "",
		Email:      "",
		Password:   string(pasword), // bcrypt hash for "password123"
		VoteStatus: false,
	}
	// Insert a dummy OTP session first
	var testOTPSession = database.OTPSession{
		ID:        "",
		VoterID:   "",
		OTP:       "",
		ExpiresAt: config.DB.NowFunc().Add(5 * time.Minute),
	}
	// Insert test voter
	config.DB.Create(&testVoter)
	config.DB.Create(testOTPSession)
}

// --- LOGIN TEST ---
func TestLoginSuccess(t *testing.T) {
	// setupTestDB()
	SetupTest(t)
	mail.LoadSMTPConfig()
	prepareTestDB()

	r := setupRouter()

	body := map[string]string{
		"username": "",
		"password": "",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	r := setupRouter()

	body := map[string]string{
		"username": "",
		"password": "",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", w.Code)
	}
}

// --- OTP VERIFY TEST ---
func TestOTPVerifySuccess(t *testing.T) {
	SetupTest(t)
	r := setupRouter()
	prepareTestDB()

	body := map[string]string{
		"voter_id":   "",
		"voter_name": "",
		"otp":        "",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestOTPVerifyInvalid(t *testing.T) {
	r := setupRouter()

	body := map[string]string{
		"voter_id": "",
		"otp":      "",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", w.Code)
	}
}

func setupFullRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/verify-otp", auth.VerifyOTP)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.SignatureMiddleware)

		r.Get("/api/rules", handler.RulesHandler)
	})

	return r
}

func TestFullFlow_VerifyOTP_GenerateJWT_ThenHitRulesPage(t *testing.T) {
	SetupTest(t)
	prepareTestDB()
	voterName := "atha"

	r := setupFullRouter()

	// 1. Simulate OTP Verification
	verifyBody := map[string]string{
		"voter_id": "",
		"otp":      "",
	}
	verifyJson, _ := json.Marshal(verifyBody)

	req1, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyJson))
	req1.Header.Set("Content-Type", "application/json")

	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("OTP verify failed, got %d", rec1.Code)
	}

	var verifyResp map[string]string
	json.NewDecoder(rec1.Body).Decode(&verifyResp)

	jwtToken := verifyResp["token"]
	if jwtToken == "" {
		t.Fatal("No token returned after OTP verification")
	}

	// 2. Simulate hitting /api/rules with JWT + Signature
	pageID := os.Getenv("PAGE_ID_RULES")
	apiKey := os.Getenv("APP_API_KEY")
	timestamp := time.Now().UTC().Format("20060102150405")

	// Build signature
	payload := voterName + "~" + pageID + "~" + timestamp + "~" + apiKey
	hash := sha256.Sum256([]byte(payload))
	signature := hex.EncodeToString(hash[:])

	req2, _ := http.NewRequest("GET", "/api/rules", nil)
	req2.Header.Set("Authorization", "Bearer "+jwtToken)
	req2.Header.Set("Signature", signature)
	req2.Header.Set("timestamp", timestamp)
	req2.Header.Set("api_key", apiKey)
	req2.Header.Set("page_id", pageID)

	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("Rules page access failed, got %d", rec2.Code)
	}

	var rulesResp map[string][]string
	json.NewDecoder(rec2.Body).Decode(&rulesResp)

	if len(rulesResp["rules"]) == 0 {
		t.Fatal("Rules content empty")
	}
}
