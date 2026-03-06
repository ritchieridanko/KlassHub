package models

import "time"

type Auth struct {
	Email             *string
	Username          *string
	Role              string
	IsVerified        bool
	PasswordChangedAt *time.Time
}
