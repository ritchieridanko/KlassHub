package dtos

import "time"

type Auth struct {
	Role              *string    `json:"role,omitempty"`
	Email             *string    `json:"email,omitempty"`
	Username          *string    `json:"username,omitempty"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Auth      *Auth      `json:"auth,omitempty"`
	AuthToken *AuthToken `json:"auth_token,omitempty"`
}
