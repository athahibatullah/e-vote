package mail

import (
	"os"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	FromName string
	ApiKey   string
}

var SmtpConfig SMTPConfig

func LoadSMTPConfig() {
	SmtpConfig = SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Pass:     os.Getenv("SMTP_PASS"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
		ApiKey:   os.Getenv("SMTP_BREVO_API_KEY"),
	}
}
