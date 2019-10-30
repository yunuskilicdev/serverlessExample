package errormessage

const (
	Ok                      = 99
	DatabaseError           = 100
	JsonParseError          = 101
	UserAlreadyExist        = 102
	UserNameOrPasswordWrong = 103
	CaptchaNeeded           = 104
	TokenIsNotValid         = 105
)

var statusText = map[int]string{
	DatabaseError:           "DATABASE_ERROR",
	JsonParseError:          "JSON_PARSE_ERROR",
	UserAlreadyExist:        "USER_ALREADY_EXIST",
	UserNameOrPasswordWrong: "USERNAME_OR_PASSWORD_WRONG",
	Ok:                      "OK",
	CaptchaNeeded:           "CAPTCHA_NEEDED",
	TokenIsNotValid:         "TOKEN_IS_NOT_VALID",
}

func StatusText(code int) string {
	return statusText[code]
}
