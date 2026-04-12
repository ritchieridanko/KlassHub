package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	SchoolUserID   *string
	Role           string
	Name           string
	Nickname       *string
	Birthplace     string
	Birthdate      *time.Time
	Sex            string
	Phone          *string
	ProfilePicture *string
	ProfileBanner  *string
	CreatedBy      *uuid.UUID
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
