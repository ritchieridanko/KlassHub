package middlewares

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/event"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

func Tracing() handlers.Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			ctx = otel.GetTextMapPropagator().Extract(
				ctx,
				event.Header(msg.Headers),
			)
			return next(ctx, msg)
		}
	}
}
