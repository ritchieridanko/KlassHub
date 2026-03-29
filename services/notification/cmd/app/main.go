package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"github.com/ritchieridanko/klasshub/services/notification/internal/di"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/handlers"
)

func main() {
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	inf, err := infra.Init(cfg)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer func(inf *infra.Infra) {
		if err := inf.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(inf)

	ctr, err := di.Init(cfg, inf)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer func(ctr *di.Container) {
		if err := ctr.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(ctr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Subscribers Start
	//
	// Topic: auth.created
	go func(ctx context.Context, s *subscriber.Subscriber, h handlers.Handler) {
		defer wg.Done()
		if err := s.Listen(ctx, h); err != nil {
			log.Println("[WARN]:", err)
		}
	}(ctx, ctr.SubscriberAC(), ctr.HandlerAC())

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	cancel()
	wg.Wait()
}
