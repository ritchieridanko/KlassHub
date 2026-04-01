package middlewares

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

func Subdomain(hostname string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		host, _, err := net.SplitHostPort(ctx.Request.Host)
		if err != nil {
			host = ctx.Request.Host
		}

		var subdomain string
		if suffix := fmt.Sprintf(".%s", hostname); strings.HasSuffix(host, suffix) {
			subdomain = strings.TrimSuffix(host, suffix)
		}

		if subdomain != constants.SubdomainAdmin && subdomain != constants.SubdomainLMS {
			ce.NewError(
				ce.CodeInvalidSubdomain,
				ce.MsgInvalidSubdomain,
				nil,
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
