package models

import "time"

type RevokeActiveSessionParams struct {
	AuthID    int64
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}

type RevokeSessionParams struct {
	RefreshToken string
	ExpiresAt    time.Time
}
