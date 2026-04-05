package middlewares

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

var roleAllowedSubdomains = map[string]string{
	constants.RoleAdministrator: constants.SubdomainAdmin,
	constants.RoleSchool:        constants.SubdomainAdmin,
	constants.RoleInstructor:    constants.SubdomainLMS,
	constants.RoleStudent:       constants.SubdomainLMS,
}

func roleAllowedSubdomain(role, subdomain string) bool {
	sd, ok := roleAllowedSubdomains[role]
	if !ok || subdomain != sd {
		return false
	}
	return true
}

func Authz(requireVerified bool, allowedSubdomains []string, allowedRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if subdomain authorization is required
		var subdomain string
		if len(allowedSubdomains) != 0 {
			subdomain = utils.CtxSubdomain(ctx.Request.Context())
			if !slices.Contains(allowedSubdomains, subdomain) {
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
		if len(allowedRoles) != 0 {
			if !slices.Contains(allowedRoles, authCtx.Role) {
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
			if subdomain != "" && !roleAllowedSubdomain(authCtx.Role, subdomain) {
				ce.NewError(
					ce.CodeNotFound,
					ce.MsgResourceNotFound,
					errors.New("subdomain unauthorized"),
					authIDField,
					schoolIDField,
					roleField,
					logger.NewField("subdomain", subdomain),
				).Bind(
					ctx,
				)

				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
