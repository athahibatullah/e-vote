package auth

import (
	"e-vote-system/internal/service"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	service.LoginService(w, r)
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	service.VerifyOTPService(w, r)
}
