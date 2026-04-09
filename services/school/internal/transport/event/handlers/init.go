package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
	"github.com/segmentio/kafka-go"
)

type Handler func(context.Context, kafka.Message) *ce.Error
type Middleware func(Handler) Handler

func NewHandler(base Handler, m ...Middleware) Handler {
	for i := len(m) - 1; i >= 0; i-- {
		base = m[i](base)
	}
	return base
}
