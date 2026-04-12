package models

import "time"

type Auth struct {
	Email             *string
	Username          *string
	Role              string
	IsVerified        bool
	SchoolExists      bool
	PasswordChangedAt *time.Time
}

type AccessToken struct {
	Token     string
	ExpiresIn int64
}

type RefreshToken struct {
	Token     string
	ExpiresIn int64
}

type AuthToken struct {
	AccessToken  *AccessToken
	RefreshToken *RefreshToken
}
