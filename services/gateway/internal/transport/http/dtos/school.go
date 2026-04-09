package dtos

import "time"

type School struct {
	NPSN           *string    `json:"npsn"`
	NPSNVerifiedAt *time.Time `json:"npsn_verified_at"`
	Name           string     `json:"name"`
	Level          string     `json:"level"`
	Ownership      string     `json:"ownership"`
	ProfilePicture *string    `json:"profile_picture"`
	ProfileBanner  *string    `json:"profile_banner"`
	Accreditation  *string    `json:"accreditation"`
	EstablishedAt  *time.Time `json:"established_at"`
	Province       string     `json:"province"`
	CityRegency    string     `json:"city_regency"`
	District       string     `json:"district"`
	Subdistrict    string     `json:"subdistrict"`
	Street         string     `json:"street"`
	Postcode       string     `json:"postcode"`
	Phone          *string    `json:"phone"`
	Email          *string    `json:"email"`
	Website        *string    `json:"website"`
	Timezone       string     `json:"timezone"`
	CreatedAt      *time.Time `json:"created_at"`
}
