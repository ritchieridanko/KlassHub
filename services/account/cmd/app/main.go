package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ritchieridanko/klasshub/services/account/configs"
	"github.com/ritchieridanko/klasshub/services/account/internal/di"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra"
	"github.com/ritchieridanko/klasshub/services/account/internal/transport/rpc/server"
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

	// Server Start
	srv := di.Init(cfg, inf).Server()
	go func(srv *server.Server) {
		if err := srv.Start(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	}(srv)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("[WARN]:", err)
	}
}
