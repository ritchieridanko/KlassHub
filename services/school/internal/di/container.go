package di

import (
	"github.com/ritchieridanko/klasshub/services/school/configs"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/school/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/school/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/school/internal/usecases"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	sdb databases.SchoolDatabase

	sr repositories.SchoolRepository

	su usecases.SchoolUsecase

	sh     *handlers.SchoolHandler
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	sdb := databases.NewSchoolDatabase(db)

	// Repositories
	sr := repositories.NewSchoolRepository(sdb)

	// Usecases
	su := usecases.NewSchoolUsecase(sr)

	// Handlers
	sh := handlers.NewSchoolHandler(su)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, l, sh)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		sdb:        sdb,
		sr:         sr,
		su:         su,
		sh:         sh,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
