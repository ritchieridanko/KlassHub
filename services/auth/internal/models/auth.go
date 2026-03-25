package models

import "time"

type Auth struct {
	ID                int64
	SchoolID          int64
	Email             *string
	Username          *string
	Password          string
	Role              string
	VerifiedAt        *time.Time
	EmailChangedAt    *time.Time
	UsernameChangedAt *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

func (a *Auth) IsVerified() bool {
	return a.VerifiedAt != nil
}

// If School ID != 0, school exists
func (a *Auth) SchoolExists() bool {
	return a.SchoolID != 0
}
