package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolUsecase interface {
	GetSchoolID(ctx context.Context, req *models.GetSchoolIDRequest) (schoolID int64, err *ce.Error)
}

type schoolUsecase struct {
	sr repositories.SchoolRepository
}

func NewSchoolUsecase(sr repositories.SchoolRepository) SchoolUsecase {
	return &schoolUsecase{sr: sr}
}

func (u *schoolUsecase) GetSchoolID(ctx context.Context, req *models.GetSchoolIDRequest) (int64, *ce.Error) {
	return u.sr.GetID(ctx, &models.GetSchoolID{AuthID: req.AuthID})
}
