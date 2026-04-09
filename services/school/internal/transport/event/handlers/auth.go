package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/events/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type AuthHandler struct {
	su usecases.SchoolUsecase
}

func NewAuthHandler(su usecases.SchoolUsecase) *AuthHandler {
	return &AuthHandler{su: su}
}

func (h *AuthHandler) OnAuthSchoolUpdateFailed(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.AuthSchoolUpdateFailed
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, ce.MsgInternalServer, err)
	}
	return h.su.OnAuthSchoolUpdateFailed(
		ctx,
		evt.GetSchoolId(),
	)
}
