package di

import (
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/clients"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/jwt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	adb databases.AuthDatabase
	sdb databases.SessionDatabase

	ar repositories.AuthRepository
	sr repositories.SessionRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	au usecases.AuthUsecase
	su usecases.SessionUsecase

	ah     *handlers.AuthHandler
	server *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(i.Database())
	tx := database.NewTransactor(i.Database())
	l := logger.NewLogger(i.Logger())

	// Services
	us := clients.NewUserService(i.UserServiceClient())

	// Databases
	adb := databases.NewAuthDatabase(db)
	sdb := databases.NewSessionDatabase(db)

	// Repositories
	ar := repositories.NewAuthRepository(adb)
	sr := repositories.NewSessionRepository(sdb)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.JWT.Duration, cfg.Auth.Duration.Session, sr, tx, v, j)
	au := usecases.NewAuthUsecase(cfg.App.Name, su, ar, us, tx, v, b)

	// Handlers
	ah := handlers.NewAuthHandler(au, l)

	// Server
	srv := server.Init(cfg.App.Name, &cfg.Server, l, ah)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		adb:        adb,
		sdb:        sdb,
		ar:         ar,
		sr:         sr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		au:         au,
		su:         su,
		ah:         ah,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
