package service

import (
	"e-vote-system/config"
	"e-vote-system/internal/mail"
	database "e-vote-system/models/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type OTPRequest struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
}

func LoginService(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var voter database.Voter
	if err := config.DB.Where("voter_name = ?", req.Username).First(&voter).Error; err != nil {
		http.Error(w, "Username not found", http.StatusUnauthorized)
		return
	}

	// Check password (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(voter.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Send OTP via email
	if err := mail.SendOTP(voter); err != nil {
		log.Println("Failed to send OTP:", err)
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent to your email"})
}

func VerifyOTPService(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VoterID string `json:"voter_id"`
		OTP     string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var session database.OTPSession
	db := config.DB
	if err := db.Where("voter_id = ? AND otp = ?", req.VoterID, req.OTP).First(&session).Error; err != nil {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		http.Error(w, "OTP expired", http.StatusUnauthorized)
		return
	}

	// ✅ Find voter's username from database
	var voter database.Voter
	if err := db.Where("id = ?", req.VoterID).First(&voter).Error; err != nil {
		http.Error(w, "Voter not found", http.StatusUnauthorized)
		return
	}

	// If you want: delete OTP after success
	db.Delete(&session)

	token, err := GenerateJWT(req.VoterID, voter.VoterName)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})

}

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY")) // Move to ENV later

func GenerateJWT(userID string, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecretKey)
}

func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
