package models

import "time"

type Auth struct {
	ID                int64
	Email             *string
	Username          *string
	Password          string
	IsSchool          bool
	LastLoginAt       *time.Time
	EmailVerifiedAt   *time.Time
	EmailChangedAt    *time.Time
	UsernameChangedAt *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time

	Role     string
	SchoolID int64
}

func (a *Auth) IsEmailVerified() bool {
	return a.EmailVerifiedAt != nil
}
