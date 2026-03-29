package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/models"
	"github.com/ritchieridanko/klasshub/services/notification/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/events/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type AuthHandler struct {
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{au: au}
}

func (h *AuthHandler) OnAuthCreated(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.AuthCreated
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, err)
	}

	eventID, err := utils.ToUUID(evt.GetEventId())
	if err != nil {
		return ce.NewError(
			ce.CodeTypeConversionFailed,
			err,
			logger.NewField("event_id", evt.GetEventId()),
		)
	}
	return h.au.OnAuthCreated(
		ctx,
		&models.AuthCreatedEventReq{
			ID:                eventID,
			Email:             evt.GetEmail(),
			VerificationToken: evt.GetVerificationToken(),
			CreatedAt:         utils.ToTime(evt.GetCreatedAt()),
		},
	)
}
