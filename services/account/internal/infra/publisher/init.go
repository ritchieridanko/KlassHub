package publisher

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(brokers string, topic string, balancer kafka.Balancer, batchSize int, batchTimeout time.Duration, l *zap.Logger) *kafka.Writer {
	w := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:      strings.Split(brokers, ","),
			Topic:        topic,
			Balancer:     balancer,
			BatchSize:    batchSize,
			BatchTimeout: batchTimeout,
			RequiredAcks: int(kafka.RequireAll),
			Async:        false,
		},
	)

	l.Sugar().Infof("[PUBLISHER] initialized (topic=%s, brokers=%s)", topic, brokers)
	return w
}
