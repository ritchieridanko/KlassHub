package di

import (
	"github.com/ritchieridanko/klasshub/services/user/configs"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/user/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/user/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/user/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/validator"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	udb databases.UserDatabase

	ur repositories.UserRepository

	validator *validator.Validator

	uu usecases.UserUsecase

	uh     *handlers.UserHandler
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	udb := databases.NewUserDatabase(db)

	// Repositories
	ur := repositories.NewUserRepository(udb)

	// Utils
	v := validator.Init()

	// Usecases
	uu := usecases.NewUserUsecase(cfg.App.Name, ur, v)

	// Handlers
	uh := handlers.NewUserHandler(uu)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, v, l, uh)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		udb:        udb,
		ur:         ur,
		validator:  v,
		uu:         uu,
		uh:         uh,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
