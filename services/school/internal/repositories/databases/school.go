package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/school/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolDatabase interface {
	Create(ctx context.Context, data *models.CreateSchoolData) (s *models.School, err *ce.Error)
	Delete(ctx context.Context, schoolID int64) (err *ce.Error)
}

type schoolDatabase struct {
	database *database.Database
}

func NewSchoolDatabase(db *database.Database) SchoolDatabase {
	return &schoolDatabase{database: db}
}

func (d *schoolDatabase) Create(ctx context.Context, data *models.CreateSchoolData) (*models.School, *ce.Error) {
	query := `
		INSERT INTO schools (
			npsn, name, level, ownership, accreditation, established_at,
			province, city_regency, district, subdistrict, street, postcode,
			phone, email, website, timezone
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16
		)
		RETURNING
			id, npsn, npsn_verified_at, name, level, ownership,
			profile_picture, profile_banner, accreditation,
			established_at, province, city_regency, district,
			subdistrict, street, postcode, phone, email, website,
			timezone, created_at
	`

	var s models.School
	err := d.database.Query(
		ctx, query,
		data.NPSN, data.Name, data.Level, data.Ownership, data.Accreditation,
		data.EstablishedAt, data.Province, data.CityRegency, data.District,
		data.Subdistrict, data.Street, data.Postcode, data.Phone, data.Email,
		data.Website, data.Timezone,
	).Scan(
		&s.ID, &s.NPSN, &s.NPSNVerifiedAt, &s.Name, &s.Level, &s.Ownership,
		&s.ProfilePicture, &s.ProfileBanner, &s.Accreditation, &s.EstablishedAt,
		&s.Province, &s.CityRegency, &s.District, &s.Subdistrict, &s.Street,
		&s.Postcode, &s.Phone, &s.Email, &s.Website, &s.Timezone, &s.CreatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create school: %w", err),
		)
	}

	return &s, nil
}

func (d *schoolDatabase) Delete(ctx context.Context, schoolID int64) *ce.Error {
	query := "DELETE FROM schools WHERE id = $1"

	err := d.database.Execute(
		ctx, query,
		schoolID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete school: %w", err)
		schoolIDField := logger.NewField("school_id", schoolID)

		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(
				ce.CodeSchoolNotFound,
				ce.MsgSchoolNotFound,
				wrappedErr,
				schoolIDField,
			)
		}
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
			schoolIDField,
		)
	}

	return nil
}
