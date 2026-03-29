package subscriber

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(appName, brokers, topic string, maxBytes int, maxWait, commitInterval time.Duration, l *zap.Logger) *kafka.Reader {
	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        strings.Split(brokers, ","),
			GroupID:        appName,
			Topic:          topic,
			MaxBytes:       maxBytes,
			MaxWait:        maxWait,
			CommitInterval: commitInterval,
		},
	)

	l.Sugar().Infof("[SUBSCRIBER] initialized (topic=%s, brokers=%s)", topic, brokers)
	return r
}
