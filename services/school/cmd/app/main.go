package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ritchieridanko/klasshub/services/school/configs"
	"github.com/ritchieridanko/klasshub/services/school/internal/di"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/school/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/school/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/shared/data"
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

	sd, err := data.LoadSchool()
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	ctr := di.Init(cfg, inf, sd)

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

	// Topic: auth.school.update.failed
	go func(ctx context.Context, s *subscriber.Subscriber, h handlers.Handler) {
		defer wg.Done()
		if err := s.Listen(ctx, h); err != nil {
			log.Println("[WARN]:", err)
		}
	}(ctx, ctr.SubscriberASUF(), ctr.HandlerASUF())

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
