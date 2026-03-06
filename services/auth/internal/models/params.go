package models

import "time"

type RevokeSessionParams struct {
	AuthID    int64
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}
