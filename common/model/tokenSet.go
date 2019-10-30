package model

import "time"

type TokenSet struct {
	AccessToken     string
	ExpireAt        time.Time
	RefreshToken    string
	RefreshExpireAt time.Time
}
