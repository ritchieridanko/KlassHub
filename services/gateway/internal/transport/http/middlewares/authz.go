package middlewares

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

func Authz(requireVerified bool, subdomains []string, allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if subdomain authorization is required
		if len(subdomains) != 0 {
			subdomain := utils.CtxSubdomain(ctx.Request.Context())
			if !slices.Contains(subdomains, subdomain) {
				ce.NewError(
					ce.CodeNotFound,
					ce.MsgResourceNotFound,
					errors.New("subdomain unauthorized"),
					logger.NewField("subdomain", subdomain),
				).Bind(
					ctx,
				)

				ctx.Abort()
				return
			}
		}

		authCtx := utils.CtxAuth(ctx.Request.Context())
		if authCtx == nil {
			ce.NewError(
				ce.CodeMissingContextValue,
				ce.MsgInternalServer,
				errors.New("auth missing from context"),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		authIDField := logger.NewField("auth_id", authCtx.AuthID)
		schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
		roleField := logger.NewField("role", authCtx.Role)

		// Check if verification is required
		if requireVerified && !authCtx.IsVerified {
			ce.NewError(
				ce.CodeAuthNotVerified,
				ce.MsgAuthNotVerified,
				nil,
				authIDField,
				schoolIDField,
				roleField,
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		// Check if role authorization is required
		if len(allowedRoles) != 0 && !slices.Contains(allowedRoles, authCtx.Role) {
			ce.NewError(
				ce.CodeUnauthorizedRole,
				ce.MsgUnauthorized,
				errors.New("role unauthorized"),
				authIDField,
				schoolIDField,
				roleField,
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
