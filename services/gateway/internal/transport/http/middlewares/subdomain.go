package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

func Subdomain(hostname string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := strings.TrimSpace(ctx.GetHeader("Origin"))
		if origin == "" {
			origin = strings.TrimSpace(ctx.GetHeader("Referer"))
		}

		u, err := url.Parse(origin)
		if err != nil || u.Hostname() == "" {
			ce.NewError(
				ce.CodeNotFound,
				ce.MsgResourceNotFound,
				fmt.Errorf("failed to parse subdomain: %w", err),
				logger.NewField("origin", origin),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		subdomain := utils.NormalizeString(strings.Split(u.Hostname(), ".")[0])
		if subdomain != constants.SubdomainAdmin && subdomain != constants.SubdomainLMS {
			ce.NewError(
				ce.CodeNotFound,
				ce.MsgResourceNotFound,
				errors.New("invalid subdomain"),
				logger.NewField("subdomain", subdomain),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(
				ctx.Request.Context(),
				constants.CtxKeySubdomain,
				subdomain,
			),
		)
		ctx.Next()
	}
}
