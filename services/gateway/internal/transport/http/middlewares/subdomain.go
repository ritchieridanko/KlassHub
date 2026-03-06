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

func Subdomain(hostname, tld string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		host, _, err := net.SplitHostPort(ctx.Request.Host)
		if err != nil {
			host = ctx.Request.Host
		}

		var subdomain string
		if suffix := fmt.Sprintf(".%s.%s", hostname, tld); strings.HasSuffix(host, suffix) {
			subdomain = strings.TrimSuffix(host, suffix)
		} else if localhost := ".localhost"; strings.HasSuffix(host, localhost) {
			subdomain = strings.TrimSuffix(host, localhost)
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
