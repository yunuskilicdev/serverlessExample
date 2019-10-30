package main

type SignupRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
