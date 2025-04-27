package handler

import (
	"e-vote-system/middleware"
	"fmt"
	"net/http"
)

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextKeyUserID).(string)
	username := r.Context().Value(middleware.ContextKeyUsername).(string)

	fmt.Fprintf(w, "Hello %s (ID: %s), you can now vote!", username, userID)
}
