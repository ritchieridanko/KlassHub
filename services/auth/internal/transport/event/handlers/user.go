package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/events/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type UserHandler struct {
	au usecases.AuthUsecase
}

func NewUserHandler(au usecases.AuthUsecase) *UserHandler {
	return &UserHandler{au: au}
}

func (h *UserHandler) OnUserCreationFailed(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.UserCreationFailed
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, ce.MsgInternalServer, err)
	}
	return h.au.OnUserCreationFailed(
		ctx,
		evt.GetAuthId(),
	)
}
