package dtos

import "time"

type CreateSchoolProfileRequest struct {
	NPSN          *string    `json:"npsn"`
	Name          string     `json:"name" binding:"required"`
	Level         string     `json:"level" binding:"required"`
	Ownership     string     `json:"ownership" binding:"required"`
	Accreditation *string    `json:"accreditation"`
	EstablishedAt *time.Time `json:"established_at" time_format:"2006-01-02"`
	Province      string     `json:"province" binding:"required"`
	CityRegency   string     `json:"city_regency" binding:"required"`
	District      string     `json:"district" binding:"required"`
	Subdistrict   string     `json:"subdistrict" binding:"required"`
	Street        string     `json:"street" binding:"required"`
	Postcode      string     `json:"postcode" binding:"required"`
	Phone         *string    `json:"phone"`
	Email         *string    `json:"email"`
	Website       *string    `json:"website"`
	Timezone      string     `json:"timezone" binding:"required"`
}

type CreateUserAccountRequest struct {
	// Auth
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password string  `json:"password" binding:"required"`
	Role     string  `json:"role" binding:"required"`

	// User
	SchoolUserID *string    `json:"school_user_id"`
	Name         string     `json:"name" binding:"required"`
	Birthplace   string     `json:"birthplace" binding:"required"`
	Birthdate    *time.Time `json:"birthdate" time_format:"2006-01-02"`
	Sex          string     `json:"sex" binding:"required"`
}
