package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/handlers"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(appName string, l *logger.Logger, ah *handlers.AuthHandler) *Router {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(appName))
	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})

	v1 := r.Group("/api/v1")

	// Auth Endpoints
	auth := v1.Group("/auth")
	{
		auth.POST("/login", ah.Login)
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
