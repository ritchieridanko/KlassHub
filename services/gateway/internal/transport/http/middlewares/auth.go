package middlewares

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/jwt"
)

func Auth(j *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if len(authorization) == 0 {
			ce.NewError(
				ce.CodeUnauthenticated,
				ce.MsgUnauthenticated,
				errors.New("access token missing from request"),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		auth := strings.Split(authorization, " ")
		if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
			ce.NewError(
				ce.CodeUnauthenticated,
				ce.MsgUnauthenticated,
				errors.New("access token is malformed"),
				logger.NewField("access_token", authorization),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		claim, err := j.Parse(auth[1])
		if err != nil {
			switch {
			case errors.Is(err, ce.ErrJWTExpired), errors.Is(err, ce.ErrJWTMalformed),
				errors.Is(err, ce.ErrInvalidJWTClaim):
				ce.NewError(
					ce.CodeUnauthenticated,
					ce.MsgUnauthenticated,
					err,
				).Bind(
					ctx,
				)
			default:
				ce.NewError(
					ce.CodeUnknown,
					ce.MsgInternalServer,
					err,
				).Bind(
					ctx,
				)
			}

			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(
				ctx.Request.Context(),
				constants.CtxKeyAuth,
				&models.AuthContext{
					AuthID:     claim.AuthID,
					SchoolID:   claim.SchoolID,
					Role:       claim.Role,
					IsVerified: claim.IsVerified,
				},
			),
		)
		ctx.Next()
	}
}
