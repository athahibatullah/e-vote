package models

import "time"

type AppConfig struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	VoterTotalCount int
	VoterHasVoted   int
	VoteDeadline    time.Time
}

func (AppConfig) TableName() string {
	return "t_app_config"
}
