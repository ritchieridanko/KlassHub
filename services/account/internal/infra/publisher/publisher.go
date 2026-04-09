package publisher

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/account/internal/utils/event"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(w *kafka.Writer) *Publisher {
	return &Publisher{writer: w}
}

func (p *Publisher) Publish(ctx context.Context, key string, msg proto.Message) error {
	value, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:     []byte(key),
			Value:   value,
			Headers: event.NewHeader(ctx),
		},
	)
}
