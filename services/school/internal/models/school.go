package models

import "time"

type School struct {
	ID             int64
	NPSN           *string
	NPSNVerifiedAt *time.Time
	Name           string
	Level          string
	Ownership      string
	ProfilePicture *string
	ProfileBanner  *string
	Accreditation  *string
	EstablishedAt  *time.Time
	Province       string
	CityRegency    string
	District       string
	Subdistrict    string
	Street         string
	Postcode       string
	Phone          *string
	Email          *string
	Website        *string
	Timezone       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
