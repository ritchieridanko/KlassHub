package main

import (
	"flag"
	"log"

	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/database"
)

func main() {
	fu := flag.Bool("up", false, "Apply all up migrations")
	fd := flag.Int("down", 0, "Apply N down migrations")
	flag.Parse()

	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	m, err := database.NewMigrator(&cfg.Database, "./migrations")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	defer func(m *database.Migrator) {
		if err := m.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(m)

	if *fu {
		if err := m.Up(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	} else if *fd >= 0 {
		if err := m.Down(*fd); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	} else {
		log.Fatalln("[FATAL]: failed to apply migrations: no action specified")
	}
}
