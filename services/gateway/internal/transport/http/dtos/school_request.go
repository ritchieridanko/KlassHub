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
