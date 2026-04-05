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
		auth.POST("/refresh", ah.RotateAuthToken)

		// Emails
		email := auth.Group("/email")
		{
			// Availability
			email.GET("/available", ah.IsEmailAvailable)

			// Verifications
			verification := email.Group("/verification")
			{
				// Resend
				verification.POST(
					"/resend",
					middlewares.Auth(j),
					middlewares.Authz(
						false,
						[]string{constants.SubdomainAdmin},
						[]string{constants.RoleSchool},
					),
					ah.ResendVerification,
				)

				// Confirm
				verification.POST(
					"/confirm",
					middlewares.Auth(j),
					middlewares.Authz(
						false,
						[]string{constants.SubdomainAdmin},
						[]string{constants.RoleSchool},
					),
					ah.VerifyEmail,
				)
			}
		}

		// Passwords
		password := auth.Group("/password")
		{
			// Change
			password.PATCH(
				"",
				middlewares.Auth(j),
				middlewares.Authz(
					true,
					constants.AllSubdomains,
					constants.AllRoles,
				),
				ah.ChangePassword,
			)
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
