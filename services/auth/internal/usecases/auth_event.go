package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"go.opentelemetry.io/otel"
)

func (u *authUsecase) OnUserCreationFailed(ctx context.Context, authID int64) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.OnUserCreationFailed")
	defer span.End()

	// Verification Token Deletion
	if err := u.tr.DeleteVerification(ctx, authID); err != nil {
		u.logger.Warn(
			ctx,
			"on user creation failed",
			logger.NewField("auth_id", authID),
			logger.NewField("error_code", err.Code()),
			logger.NewField("error", err.Unwrap()),
		)
	}

	// Auth Deletion
	if err := u.ar.Delete(ctx, authID); err != nil && err.Code() != ce.CodeAuthNotFound {
		return err
	}
	return nil
}
