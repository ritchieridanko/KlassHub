package usecases

import (
	"context"
	"time"

	"github.com/ritchieridanko/klasshub/services/notification/internal/channels"
	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/models"
	"github.com/ritchieridanko/klasshub/services/notification/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"go.opentelemetry.io/otel"
)

type AuthUsecase interface {
	OnAuthCreated(ctx context.Context, req *models.ACEventReq) (err *ce.Error)
	OnAuthVerificationRequested(ctx context.Context, req *models.AVREventReq) (err *ce.Error)
}

type authUsecase struct {
	appName     string
	mailTimeout time.Duration
	ec          channels.EmailChannel
	er          repositories.EventRepository
	logger      *logger.Logger
}

func NewAuthUsecase(appName string, mailTimeout time.Duration, ec channels.EmailChannel, er repositories.EventRepository, l *logger.Logger) AuthUsecase {
	return &authUsecase{
		appName:     appName,
		mailTimeout: mailTimeout,
		ec:          ec,
		er:          er,
		logger:      l,
	}
}

func (u *authUsecase) OnAuthCreated(ctx context.Context, req *models.ACEventReq) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.OnAuthCreated")
	defer span.End()

	eventIDField := logger.NewField("event_id", req.ID.String())

	// Event Record Fetching
	evt, err := u.er.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// Idempotency Check
	if evt == nil {
		// Event Record Creation
		rm, err := utils.ToRawMessage(req)
		if err != nil {
			return ce.NewError(ce.CodeJSONRawEncodingFailed, err, eventIDField)
		}

		createErr := u.er.Create(
			ctx,
			&models.CreateEventData{
				ID:      req.ID,
				Topic:   constants.EventTopicAC,
				Payload: rm,
			},
		)
		if createErr != nil {
			return createErr
		}
	}
	if evt != nil {
		// Completion Status Check
		if evt.CompletedAt != nil {
			return nil
		}

		// Mailing Timeout Check
		if time.Since(evt.LastProcessedAt).Seconds() < u.mailTimeout.Seconds() {
			return ce.NewError(ce.CodeEventOnProcess, ce.ErrEventOnProcess, eventIDField)
		}

		// Event Record Update
		if err := u.er.SetLastProcessed(ctx, evt.ID); err != nil {
			return err
		}
	}

	// Email Delivery
	err = u.ec.SendWelcome(
		ctx,
		&models.WelcomeEmailMsg{
			Recipient:         req.Email,
			VerificationToken: req.VerificationToken,
		},
	)
	if err != nil {
		return err.Append(eventIDField)
	}

	u.logger.Info(ctx, "EMAIL SENT", eventIDField)

	// Event Record Update
	return u.er.SetCompleted(ctx, req.ID)
}

func (u *authUsecase) OnAuthVerificationRequested(ctx context.Context, req *models.AVREventReq) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.OnAuthVerificationRequested")
	defer span.End()

	eventIDField := logger.NewField("event_id", req.ID.String())

	// Event Record Fetching
	evt, err := u.er.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// Idempotency Check
	if evt == nil {
		// Event Record Creation
		rm, err := utils.ToRawMessage(req)
		if err != nil {
			return ce.NewError(ce.CodeJSONRawEncodingFailed, err, eventIDField)
		}

		createErr := u.er.Create(
			ctx,
			&models.CreateEventData{
				ID:      req.ID,
				Topic:   constants.EventTopicAVR,
				Payload: rm,
			},
		)
		if createErr != nil {
			return createErr
		}
	}
	if evt != nil {
		// Completion Status Check
		if evt.CompletedAt != nil {
			return nil
		}

		// Mailing Timeout Check
		if time.Since(evt.LastProcessedAt).Seconds() < u.mailTimeout.Seconds() {
			return ce.NewError(ce.CodeEventOnProcess, ce.ErrEventOnProcess, eventIDField)
		}

		// Event Record Update
		if err := u.er.SetLastProcessed(ctx, evt.ID); err != nil {
			return err
		}
	}

	// Email Delivery
	err = u.ec.SendVerification(
		ctx,
		&models.VerificationEmailMsg{
			Recipient:         req.Email,
			VerificationToken: req.VerificationToken,
		},
	)
	if err != nil {
		return err.Append(eventIDField)
	}

	u.logger.Info(ctx, "EMAIL SENT", eventIDField)

	// Event Record Update
	return u.er.SetCompleted(ctx, req.ID)
}
