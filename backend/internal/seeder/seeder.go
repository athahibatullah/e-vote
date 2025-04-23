package seeder

import (
	models "e-vote-system/models/database"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(db *gorm.DB) {
	// Auto-migrate tables (you can remove this if migrations are separate)
	db.AutoMigrate(&models.AppConfig{}, &models.Candidate{}, &models.Voter{})

	// Seed app config
	appConfigs := []models.AppConfig{
		{
			VoterTotalCount: 100,
			VoterHasVoted:   0,
			VoteDeadline:    time.Date(2025, 5, 1, 23, 59, 0, 0, time.UTC)},
	}
	for _, c := range appConfigs {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&c)
	}

	// Seed candidates
	candidates := []models.Candidate{
		{
			ID:            uuid.NewString(),
			CandidateName: "John Doe",
			VisionMission: "Peace and Progress",
		},
		{
			ID:            uuid.NewString(),
			CandidateName: "Jane Smith",
			VisionMission: "Innovation and Equality",
		},
	}

	for _, c := range candidates {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&c)
	}

	// Seed voters (optional)
	voters := []models.Voter{
		{
			ID:         uuid.NewString(),
			VoterName:  "Alice",
			Password:   hashPassword(""),
			VoteStatus: false,
		},
		{
			ID:         uuid.NewString(),
			VoterName:  "Bob",
			Password:   hashPassword(""),
			VoteStatus: false,
		},
	}

	for _, v := range voters {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&v)
	}

	log.Println("✅ Seed data inserted.")
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}
	return string(hashed)
}
