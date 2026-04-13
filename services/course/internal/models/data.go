package models

import "github.com/google/uuid"

type CreateCourseData struct {
	ID             uuid.UUID
	SchoolID       int64
	SchoolCourseID *string
	Name           string
	Description    *string
	CoursePicture  *string
}
