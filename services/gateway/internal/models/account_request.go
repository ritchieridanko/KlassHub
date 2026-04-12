package models

import "time"

type CreateSchoolProfileReq struct {
	NPSN          *string
	Name          string
	Level         string
	Ownership     string
	Accreditation *string
	EstablishedAt *time.Time
	Province      string
	CityRegency   string
	District      string
	Subdistrict   string
	Street        string
	Postcode      string
	Phone         *string
	Email         *string
	Website       *string
	Timezone      string

	// For Auth Token Refresh
	RefreshToken string
}

type CreateUserAccountReq struct {
	// Auth
	Email    *string
	Username *string
	Password string
	Role     string

	// User
	SchoolUserID *string
	Name         string
	Birthplace   string
	Birthdate    *time.Time
	Sex          string
}
