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

func Init(cfg *configs.Client, appName string, j *jwt.JWT, l *logger.Logger, auh *handlers.AuthHandler, ach *handlers.AccountHandler) *Router {
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
		auth.POST("/login", auh.Login)
		auth.POST("/logout", middlewares.Auth(j), auh.Logout)
		auth.POST("/register", auh.CreateSchoolAuth)
		auth.POST("/refresh", auh.RotateAuthToken)

		// Usernames
		username := auth.Group("/username")
		{
			// Availability
			username.GET(
				"/available",
				middlewares.Auth(j),
				middlewares.Authz(
					false,
					[]string{constants.SubdomainAdmin},
					constants.AdminRoles,
				),
				auh.IsUsernameAvailable,
			)
		}

		// Emails
		email := auth.Group("/email")
		{
			// Availability
			email.GET("/available", auh.IsEmailAvailable)

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
					auh.ResendVerification,
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
					auh.VerifyEmail,
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
				auh.ChangePassword,
			)
		}
	}

	// SCHOOL ENDPOINTS
	school := v1.Group("/schools")
	{
		// Create
		school.POST(
			"",
			middlewares.Auth(j),
			middlewares.Authz(
				true,
				[]string{constants.SubdomainAdmin},
				[]string{constants.RoleSchool},
			),
			ach.CreateSchoolProfile,
		)
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
