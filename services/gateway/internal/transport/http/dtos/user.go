package dtos

import "time"

type User struct {
	ID             string     `json:"id"`
	SchoolUserID   *string    `json:"school_user_id"`
	Role           string     `json:"role"`
	Name           string     `json:"name"`
	Nickname       *string    `json:"nickname"`
	Birthplace     string     `json:"birthplace"`
	Birthdate      *time.Time `json:"birthdate"`
	Sex            string     `json:"sex"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	ProfileBanner  *string    `json:"profile_banner"`
}

type UserAdmin struct {
	ID             string     `json:"id"`
	SchoolUserID   *string    `json:"school_user_id"`
	Role           string     `json:"role"`
	Name           string     `json:"name"`
	Birthplace     string     `json:"birthplace"`
	Birthdate      *time.Time `json:"birthdate"`
	Sex            string     `json:"sex"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	CreatedBy      *string    `json:"created_by"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
