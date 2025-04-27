package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string
	Username  string
	Password  string
	OTP       string
	OTPExpiry time.Time
}

func FindUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func SaveOTP(db *gorm.DB, userID string, otp string, expiry time.Time) error {
	return db.Model(&User{}).Where("id = ?", userID).Updates(User{OTP: otp, OTPExpiry: expiry}).Error
}
