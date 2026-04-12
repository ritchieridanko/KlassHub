package subscriber

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	processTimeout time.Duration
	reader         *kafka.Reader
	logger         *logger.Logger
}

func NewSubscriber(processTimeout time.Duration, r *kafka.Reader, l *logger.Logger) *Subscriber {
	return &Subscriber{
		processTimeout: processTimeout,
		reader:         r,
		logger:         l,
	}
}

func (s *Subscriber) Listen(ctx context.Context, handler func(context.Context, kafka.Message) *ce.Error) error {
	topic := s.reader.Config().Topic
	partition := s.reader.Config().Partition
	topicField := logger.NewField("topic", topic)
	partitionField := logger.NewField("partition", partition)

	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf(
					"failed to fetch message (topic=%s, partition=%d): %w",
					topic, partition, err,
				)
			}

			s.logger.Error(
				ctx,
				"EVENT FETCHING ERROR",
				topicField,
				partitionField,
				logger.NewField("error_code", ce.CodeEventFetchingFailed),
				logger.NewField("error", err.Error()),
			)
			continue
		}

		c, cancel := context.WithTimeout(ctx, s.processTimeout)
		handleErr := handler(c, msg)
		if cancel(); handleErr != nil {
			continue
		}
		if err := s.reader.CommitMessages(ctx, msg); err != nil {
			s.logger.Error(
				ctx,
				"EVENT COMMITTING ERROR",
				topicField,
				partitionField,
				logger.NewField("error_code", ce.CodeEventCommittingFailed),
				logger.NewField("error", err.Error()),
			)
		}
	}
}
