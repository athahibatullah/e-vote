package main

import (
	"e-vote-system/config"
	"e-vote-system/internal/seeder"
	"fmt"
	"net/http"
	"os"
)

func main() {
	config.ConnectDatabase()
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		db := config.DB
		seeder.Seed(db)
		fmt.Println("✅ Seeding completed")
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "E-Vote backend is alive!")
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
