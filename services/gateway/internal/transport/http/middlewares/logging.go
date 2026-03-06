package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

func Logging(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		fields := []logger.Field{
			logger.NewField("method", ctx.Request.Method),
			logger.NewField("path", ctx.FullPath()),
			logger.NewField("latency", time.Since(start).String()),
		}

		errs := ctx.Errors
		if len(errs) == 0 {
			fields = append(fields, logger.NewField("status", ctx.Writer.Status()))
			l.Info(ctx.Request.Context(), "REQUEST OK", fields...)
			return
		}

		var e *ce.Error
		if errors.As(errs[0].Err, &e) {
			status := e.ToHTTPStatus()

			fields = append(fields, e.Fields()...)
			fields = append(
				fields,
				logger.NewField("status", status),
				logger.NewField("error_code", e.Code()),
				logger.NewField("error", e.Error()),
			)

			if status >= http.StatusInternalServerError {
				l.Error(ctx.Request.Context(), "SYSTEM ERROR", fields...)
			} else {
				l.Warn(ctx.Request.Context(), "REQUEST ERROR", fields...)
			}

			utils.SetResponse[any](ctx, status, e.Message(), nil)
			return
		}

		// Fallback
		fields = append(
			fields,
			logger.NewField("status", http.StatusInternalServerError),
			logger.NewField("error_code", ce.CodeUnknown),
			logger.NewField("error", errs[0].Err.Error()),
		)

		l.Error(ctx.Request.Context(), "UNKNOWN ERROR", fields...)
		utils.SetResponse[any](ctx, http.StatusInternalServerError, ce.MsgInternalServer, nil)
	}
}
