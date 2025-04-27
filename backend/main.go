package main

import (
	"e-vote-system/config"
	"e-vote-system/internal/auth"
	"e-vote-system/internal/handler"
	"e-vote-system/internal/mail"
	"e-vote-system/internal/seeder"
	jwtmiddleware "e-vote-system/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.ConnectDatabase()
	mail.LoadSMTPConfig()

	// Run seeder if "seed" argument is passed
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		db := config.DB
		seeder.Seed(db)
		fmt.Println("✅ Seeding completed")
		return
	}

	// Create new router
	r := chi.NewRouter()

	// Basic middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "E-Vote backend is alive!")
	})

	// Auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", auth.Login)
		r.Post("/verify-otp", auth.VerifyOTP)
	})

	// (later) Candidate voting routes, user profile routes, etc.
	r.Group(func(r chi.Router) {
		r.Use(jwtmiddleware.JWTMiddleware)
		r.Use(jwtmiddleware.SignatureMiddleware)

		r.Get("/api/vote", handler.VoteHandler) // <- protected route
		r.Get("/api/rules", handler.RulesHandler)
		r.Get("/api/vision-mission", handler.VisionMissionHandler)
		// r.Get("/api/vote/show", handler.VoteShowHandler)
		// r.Post("/api/vote/submit", handler.VoteSubmitHandler)
	})

	// Start server
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
