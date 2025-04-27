package middleware

import (
	"context"
	"e-vote-system/internal/service"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "userID"
	ContextKeyUsername contextKey = "username"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := fields[1]

		claims, err := service.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Inject user info into context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims["sub"])
		ctx = context.WithValue(ctx, ContextKeyUsername, claims["username"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
