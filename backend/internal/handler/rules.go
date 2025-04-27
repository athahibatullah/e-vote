package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"e-vote-system/middleware"
)

func RulesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Validate JWT token (already done by middleware)

	// 2. Validate Signature
	username := r.Context().Value(middleware.ContextKeyUsername).(string)

	signature := r.Header.Get("Signature")
	timestamp := r.Header.Get("timestamp")
	apiKey := r.Header.Get("api_key")
	pageID := r.Header.Get("page_id") // or pass it through header if you prefer

	if signature == "" || timestamp == "" || apiKey == "" || pageID == "" {
		http.Error(w, "Missing required headers or params", http.StatusBadRequest)
		return
	}

	// Check timestamp within ±1 minute
	tsInt, err := time.Parse("20060102150405", timestamp)
	if err != nil {
		http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if now.Sub(tsInt) > 1*time.Minute || tsInt.Sub(now) > 1*time.Minute {
		http.Error(w, "Request expired", http.StatusUnauthorized)
		return
	}

	// Rebuild the payload
	payload := username + "~" + pageID + "~" + timestamp + "~" + apiKey
	hash := sha256.Sum256([]byte(payload))
	expectedSignature := hex.EncodeToString(hash[:])

	if signature != expectedSignature {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// 3. Serve Rules Content
	rawRules := os.Getenv("VOTING_RULES")

	rules := strings.Split(rawRules, ".")
	var cleanedRules []string
	for _, rule := range rules {
		trimmed := strings.TrimSpace(rule)
		if trimmed != "" {
			cleanedRules = append(cleanedRules, trimmed+".")
		}
	}

	resp := map[string][]string{
		"rules": cleanedRules,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
