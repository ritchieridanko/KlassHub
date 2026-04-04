package models

import (
	"time"

	"github.com/google/uuid"
)

type ACEventReq struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	VerificationToken string     `json:"verification_token"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
}

type AVREventReq struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	VerificationToken string     `json:"verification_token"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
}
