package models

import "time"

type OTPSession struct {
	ID        string    `gorm:"primaryKey"`
	VoterID   string    `gorm:"not null"`
	OTP       string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (OTPSession) TableName() string {
	return "t_otp_session"
}
