package test

import (
	"e-vote-system/config"
	models "e-vote-system/models/database"
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func SetupTest(t *testing.T) {
	// Load environment variables for testing
	err := godotenv.Load("../../../.env.local")
	if err != nil {
		log.Fatalf("Failed to load .env.local: %v", err)
	}

	// Connect to the database
	config.ConnectDatabase()

	// Auto-migrate tables needed
	err = config.DB.AutoMigrate(
		&models.Voter{},
		&models.OTPSession{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Truncate tables before each test
	CleanDatabase()
}

func CleanDatabase() {
	models := []string{"t_otp_session", "t_voter_detail"}

	for _, table := range models {
		err := config.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)).Error
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}
