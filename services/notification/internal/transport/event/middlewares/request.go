package middlewares

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Request() handlers.Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			ctx = context.WithValue(ctx, constants.CtxKeyRequestID, utils.MsgRequestID(msg))
			return next(ctx, msg)
		}
	}
}
