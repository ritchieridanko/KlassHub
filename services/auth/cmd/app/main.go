package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/di"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/server"
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

	ctr := di.Init(cfg, inf)

	// Server Start
	srv := ctr.Server()
	go func(srv *server.Server) {
		if err := srv.Start(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	}(srv)

	// Subscribers Start
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Topic: user.creation.failed
	go func(ctx context.Context, s *subscriber.Subscriber, h handlers.Handler) {
		defer wg.Done()
		if err := s.Listen(ctx, h); err != nil {
			log.Println("[WARN]:", err)
		}
	}(ctx, ctr.SubscriberUCF(), ctr.HandlerUCF())

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	cancel()
	wg.Wait()

	sdCtx, sdCancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer sdCancel()

	if err := srv.Shutdown(sdCtx); err != nil {
		log.Println("[WARN]:", err)
	}
}
