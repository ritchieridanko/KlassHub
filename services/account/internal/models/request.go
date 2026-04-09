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

type UpdateSchoolReq struct {
	SchoolID     int64
	RefreshToken string
}
