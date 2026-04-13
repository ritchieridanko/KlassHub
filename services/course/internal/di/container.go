package di

import (
	"github.com/ritchieridanko/klasshub/services/course/configs"
	"github.com/ritchieridanko/klasshub/services/course/internal/clients"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/course/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/course/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/course/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/course/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/course/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/validator"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	sc clients.SchoolClient

	cdb databases.CourseDatabase

	cr repositories.CourseRepository

	validator *validator.Validator

	cu usecases.CourseUsecase

	ch     *handlers.CourseHandler
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Clients
	sc := clients.NewSchoolClient(inf.SchoolService())

	// Databases
	cdb := databases.NewCourseDatabase(db)

	// Repositories
	cr := repositories.NewCourseRepository(cdb)

	// Utils
	v := validator.Init()

	// Usecases
	cu := usecases.NewCourseUsecase(cfg.App.Name, sc, cr, v)

	// Handlers
	ch := handlers.NewCourseHandler(cu)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, v, l, ch)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		sc:         sc,
		cdb:        cdb,
		cr:         cr,
		validator:  v,
		cu:         cu,
		ch:         ch,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
