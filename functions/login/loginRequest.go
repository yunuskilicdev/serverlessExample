package main

type LoginRequest struct {
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	CaptchaId       string
	CaptchaResponse string
}
