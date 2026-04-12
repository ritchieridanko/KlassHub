package middlewares

import (
	"context"
	"time"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Logging(l *logger.Logger) handlers.Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			start := time.Now()
			err := next(ctx, msg)

			fields := []logger.Field{
				logger.NewField("topic", msg.Topic),
				logger.NewField("partition", msg.Partition),
				logger.NewField("offset", msg.Offset),
				logger.NewField("key", string(msg.Key)),
				logger.NewField("latency", time.Since(start).String()),
			}
			if err != nil {
				fields = append(fields, err.Fields()...)
				fields = append(
					fields,
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Error()),
				)

				l.Error(ctx, "EVENT PROCESSING ERROR", fields...)
				return err
			}

			l.Info(ctx, "EVENT PROCESSING OK", fields...)
			return nil
		}
	}
}
