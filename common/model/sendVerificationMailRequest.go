package model

type SendVerificationMailRequest struct {
	UserId uint
	Email  string
	Token  string
}
