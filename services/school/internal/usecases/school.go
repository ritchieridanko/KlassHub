package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolUsecase interface {
	GetID(ctx context.Context, authID int64) (schoolID int64, err *ce.Error)
}

type schoolUsecase struct {
	sr repositories.SchoolRepository
}

func NewSchoolUsecase(sr repositories.SchoolRepository) SchoolUsecase {
	return &schoolUsecase{sr: sr}
}

func (u *schoolUsecase) GetID(ctx context.Context, authID int64) (int64, *ce.Error) {
	return u.sr.GetID(ctx, authID)
}
