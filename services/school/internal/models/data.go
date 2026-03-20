package models

import "time"

type CreateSchoolData struct {
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
