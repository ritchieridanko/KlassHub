package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/dtos"
)

func SetResponse[T any](ctx *gin.Context, status int, message string, data T) {
	ctx.JSON(
		status,
		dtos.Response[T]{
			Status:  status,
			Message: message,
			Data:    data,
			Metadata: &dtos.ResponseMetadata{
				RequestID: CtxRequestID(ctx.Request.Context()),
				Timestamp: time.Now().UTC(),
			},
		},
	)
}
