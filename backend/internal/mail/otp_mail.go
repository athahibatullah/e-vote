package mail

import (
	"crypto/rand"
	"e-vote-system/config"
	database "e-vote-system/models/database"
	"fmt"
	"math/big"
	"net/smtp"
	"time"

	"github.com/google/uuid"
)

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp += string(digits[num.Int64()])
	}
	return otp, nil
}

func SendOTP(voter database.Voter) error {
	// Generate random OTP
	otp, err := generateOTP(6)
	if err != nil {
		return err
	}

	// Save OTP into database
	db := config.DB
	session := database.OTPSession{
		ID:        uuid.NewString(),
		VoterID:   voter.ID,
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5 min expiry
	}
	if err := db.Create(&session).Error; err != nil {
		return err
	}

	// Send OTP Email
	auth := smtp.PlainAuth("", SmtpConfig.User, SmtpConfig.Pass, SmtpConfig.Host)

	to := []string{voter.Email} // make sure voter.VoterName is an email
	msg := []byte(fmt.Sprintf(
		"Subject: Your OTP Code\n"+
			"From: %s\n"+
			"To: %s\n\n"+
			"Your OTP is: %s\n"+
			"This OTP will expire in 5 minutes.\n",
		SmtpConfig.FromName,
		voter.VoterName,
		otp,
	))

	err = smtp.SendMail(SmtpConfig.Host+":"+SmtpConfig.Port, auth, SmtpConfig.User, to, msg)
	if err != nil {
		return err
	}

	// msg := gomail.NewMessage()
	// msg.SetHeader("From", SmtpConfig.User)
	// msg.SetHeader("To", voter.Email)
	// msg.SetHeader("Subject", "EVoteNation OTP Code")
	// msg.SetBody("text/plain", fmt.Sprintf("Your OTP is: %s\nThis OTP will expire in 5 minutes.", otp))

	// n := gomail.NewDialer(SmtpConfig.Host, 587, SmtpConfig.User, SmtpConfig.ApiKey)
	// if err := n.DialAndSend(msg); err != nil {
	// 	panic(err)
	// }
	fmt.Println("✅ OTP sent and stored for", voter.VoterName)
	return nil
}
