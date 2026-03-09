package models

import "time"

type CreateAuthData struct {
	SchoolID   int64
	Email      *string
	Username   *string
	Password   string
	Role       string
	VerifiedAt *time.Time
}

type CreateSessionData struct {
	ParentID     *int64
	AuthID       int64
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
}
