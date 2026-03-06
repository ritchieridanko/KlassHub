package dtos

import "time"

type Auth struct {
	Email             *string    `json:"email,omitempty"`
	Username          *string    `json:"username,omitempty"`
	Role              string     `json:"role"`
	IsVerified        bool       `json:"is_verified"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
}
