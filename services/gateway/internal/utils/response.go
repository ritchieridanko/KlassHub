package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Status   int               `json:"status"`
	Message  string            `json:"message"`
	Data     T                 `json:"data,omitempty"`
	Metadata *ResponseMetadata `json:"metadata,omitempty"`
}

type ResponseMetadata struct {
	RequestID string    `json:"request_id"`
	Page      *int      `json:"page,omitempty"`
	PageSize  *int      `json:"page_size,omitempty"`
	Total     *int      `json:"total,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func SetResponse[T any](ctx *gin.Context, status int, message string, data T) {
	ctx.JSON(
		status,
		Response[T]{
			Status:  status,
			Message: message,
			Data:    data,
			Metadata: &ResponseMetadata{
				RequestID: CtxRequestID(ctx.Request.Context()),
				Timestamp: time.Now().UTC(),
			},
		},
	)
}
