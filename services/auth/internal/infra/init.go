package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	cache    *redis.Client
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer
	acp      *kafka.Writer
	avrp     *kafka.Writer
	ucfs     *kafka.Reader
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	cc, err := cache.Init(&cfg.Cache, l)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Publishers
	acp := publisher.Init(
		cfg.Broker.Brokers,
		cfg.Broker.Publisher.AC.Name,
		&kafka.Murmur2Balancer{
			Consistent: true,
		},
		cfg.Broker.Publisher.AC.BatchSize,
		cfg.Broker.Publisher.AC.BatchTimeout,
		l,
	)
	avrp := publisher.Init(
		cfg.Broker.Brokers,
		cfg.Broker.Publisher.AVR.Name,
		&kafka.Murmur2Balancer{
			Consistent: true,
		},
		cfg.Broker.Publisher.AVR.BatchSize,
		cfg.Broker.Publisher.AVR.BatchTimeout,
		l,
	)

	// Subscribers
	ucfs := subscriber.Init(
		cfg.App.Name,
		cfg.Broker.Brokers,
		cfg.Broker.Subscriber.UCF.Name,
		cfg.Broker.Subscriber.UCF.MaxBytes,
		cfg.Broker.Subscriber.UCF.MaxWait,
		cfg.Broker.Subscriber.UCF.CommitInterval,
		l,
	)

	return &Infra{
		config:   cfg,
		cache:    cc,
		database: db,
		logger:   l,
		tracer:   t,
		acp:      acp,
		avrp:     avrp,
		ucfs:     ucfs,
	}, nil
}

func (i *Infra) Cache() *redis.Client {
	return i.cache
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) PublisherAC() *kafka.Writer {
	return i.acp
}

func (i *Infra) PublisherAVR() *kafka.Writer {
	return i.avrp
}

func (i *Infra) SubscriberUCF() *kafka.Reader {
	return i.ucfs
}

func (i *Infra) Close() error {
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache: %w", err)
	}
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.acp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicAC, err)
	}
	if err := i.avrp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicAVR, err)
	}
	if err := i.ucfs.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber (topic: %s): %w", constants.EventTopicUCF, err)
	}

	i.database.Close()
	return nil
}
