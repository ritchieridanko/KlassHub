package models

import "time"

type CreateSchoolReq struct {
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
}

type CreateSchoolProfileReq struct {
	// School
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

	// Auth
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

type CreateUserAuthReq struct {
	Email    *string
	Username *string
	Password string
	Role     string
}

type CreateUserReq struct {
	AuthID       int64
	SchoolID     int64
	SchoolUserID *string
	Role         string
	Name         string
	Birthplace   string
	Birthdate    *time.Time
	Sex          string
}

type UpdateSchoolReq struct {
	SchoolID     int64
	RefreshToken string
}
