package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserData struct {
	ID           uuid.UUID
	AuthID       int64
	SchoolID     int64
	SchoolUserID *string
	Role         string
	Name         string
	Birthplace   string
	Birthdate    time.Time
	Sex          string
	CreatedBy    *uuid.UUID
}
