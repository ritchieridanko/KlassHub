package dtos

import "time"

type Auth struct {
	Email             *string    `json:"email,omitempty"`
	Username          *string    `json:"username,omitempty"`
	Role              string     `json:"role"`
	IsVerified        bool       `json:"is_verified"`
	SchoolExists      bool       `json:"school_exists"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
}

type AccessToken struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
