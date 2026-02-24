package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	AuthID         int64
	SchoolID       int64
	SchoolUserID   *string
	Role           string
	Name           string
	Nickname       *string
	Birthplace     *string
	Birthdate      *time.Time
	Sex            *string
	Phone          *string
	ProfilePicture *string
	ProfileBanner  *string
	CreatedBy      *uuid.UUID
	CreatedByName  *string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
}

type GetUser struct {
	AuthID   int64
	SchoolID int64
}
