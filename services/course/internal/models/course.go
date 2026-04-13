package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID             uuid.UUID
	SchoolID       int64
	SchoolCourseID *string
	Name           string
	Description    *string
	CoursePicture  *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
