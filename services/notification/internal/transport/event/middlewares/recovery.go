package middlewares

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Recovery(l *logger.Logger) handlers.Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return func(ctx context.Context, msg kafka.Message) (err *ce.Error) {
			defer func() {
				if r := recover(); r != nil {
					l.Error(
						ctx,
						"PANIC RECOVERED",
						logger.NewField("topic", msg.Topic),
						logger.NewField("partition", msg.Partition),
						logger.NewField("offset", msg.Offset),
						logger.NewField("key", string(msg.Key)),
						logger.NewField("panic", fmt.Sprintf("%v", r)),
						logger.NewField("stack_trace", debug.Stack()),
					)
					err = ce.NewError(ce.CodePanicOccurred, fmt.Errorf("%v", r))
				}
			}()

			return next(ctx, msg)
		}
	}
}
