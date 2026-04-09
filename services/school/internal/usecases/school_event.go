package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
	"go.opentelemetry.io/otel"
)

func (u *schoolUsecase) OnAuthSchoolUpdateFailed(ctx context.Context, schoolID int64) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "school.usecase.OnAuthSchoolUpdateFailed")
	defer span.End()

	// School Deletion
	if err := u.sr.Delete(ctx, schoolID); err != nil && err.Code() != ce.CodeSchoolNotFound {
		return err
	}
	return nil
}
