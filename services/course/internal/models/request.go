package models

type CreateCourseReq struct {
	SchoolCourseID *string
	Name           string
	Description    *string
	CoursePicture  *string
}
