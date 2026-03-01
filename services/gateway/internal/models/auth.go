package models

import "time"

type Auth struct {
	Role              *string
	Email             *string
	Username          *string
	EmailVerifiedAt   *time.Time
	PasswordChangedAt *time.Time
}
