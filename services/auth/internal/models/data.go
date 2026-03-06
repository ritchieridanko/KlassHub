package models

import "time"

type CreateSessionData struct {
	ParentID     *int64
	AuthID       int64
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
}
