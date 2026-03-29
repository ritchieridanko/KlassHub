package event

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

type Header []kafka.Header

func NewHeader(ctx context.Context) Header {
	headers := Header{
		kafka.Header{
			Key:   constants.EventHeaderKeyRequestID,
			Value: []byte(utils.CtxRequestID(ctx)),
		},
		kafka.Header{
			Key:   "content-type",
			Value: []byte("application/x-protobuf"),
		},
	}

	otel.GetTextMapPropagator().Inject(ctx, &headers)
	return headers
}

func (h Header) Get(key string) string {
	for _, header := range h {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

func (h Header) Set(key, value string) {
	h = append(
		h,
		kafka.Header{
			Key:   key,
			Value: []byte(value),
		},
	)
}

func (h Header) Keys() []string {
	keys := make([]string, len(h))
	for i, header := range h {
		keys[i] = header.Key
	}
	return keys
}
