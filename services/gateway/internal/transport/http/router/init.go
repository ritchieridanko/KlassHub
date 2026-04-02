package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/configs"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/middlewares"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/jwt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(cfg *configs.Client, appName string, j *jwt.JWT, l *logger.Logger, ah *handlers.AuthHandler) *Router {
	r := gin.New()
	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})

	v1 := r.Group(
		"/v1",
		middlewares.Request(l),
		middlewares.Recovery(l),
		otelgin.Middleware(appName),
		middlewares.Logging(l),
		middlewares.Subdomain(cfg.Host),
	)

	// AUTH ENDPOINTS
	auth := v1.Group("/auth")
	{
		// Authentications
		auth.POST("/login", ah.Login)
		auth.POST("/logout", middlewares.Auth(j), ah.Logout)
		auth.POST("/register", ah.CreateSchoolAuth)

		email := auth.Group("/email")
		{
			// Email Verification
			email.POST(
				"/verification/confirm",
				middlewares.Auth(j),
				middlewares.Authz(
					false,
					[]string{constants.SubdomainAdmin},
					constants.RoleSchool,
				),
				ah.VerifyEmail,
			)
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
