package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/http"
	"os"
	"time"
)

func SignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(ContextKeyUsername).(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized - username missing", http.StatusUnauthorized)
			return
		}

		signature := r.Header.Get("Signature")
		timestampStr := r.Header.Get("timestamp")
		apiKey := r.Header.Get("api_key")
		pageID := r.Header.Get("page_id")

		if signature == "" || timestampStr == "" || apiKey == "" || pageID == "" {
			http.Error(w, "Missing signature-related headers or params", http.StatusBadRequest)
			return
		}

		// Parse timestamp
		tsParsed, err := time.ParseInLocation("20060102150405", r.Header.Get("timestamp"), time.UTC)
		if err != nil {
			http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()

		// calculate delta in seconds
		diffSeconds := math.Abs(now.Sub(tsParsed).Seconds())
		if diffSeconds > 60 { // more than 60 seconds
			http.Error(w, "Request expired", http.StatusUnauthorized)
			return
		}

		// Rebuild signature
		endpoint := r.URL.Path
		var pageIDConstruct string
		switch endpoint {
		case "/api/rules":
			pageIDConstruct = os.Getenv("PAGE_ID_RULES")
		case "/api/vision-mission":
			pageIDConstruct = os.Getenv("PAGE_ID_VISION_MISSION")
		case "/api/vote/show":
			pageIDConstruct = os.Getenv("PAGE_ID_VOTE_SHOW")
		default:
			http.Error(w, "Invalid URL", http.StatusUnauthorized)
			return
		}
		if pageID != pageIDConstruct {
			http.Error(w, "Invalid header", http.StatusUnauthorized)
			return
		}
		apiKeyConstruct := os.Getenv("APP_API_KEY")
		if apiKey != apiKeyConstruct {
			http.Error(w, "Invalid header", http.StatusUnauthorized)
			return
		}
		payload := username + "~" + pageIDConstruct + "~" + timestampStr + "~" + apiKeyConstruct
		hash := sha256.Sum256([]byte(payload))
		expectedSignature := hex.EncodeToString(hash[:])

		if signature != expectedSignature {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
