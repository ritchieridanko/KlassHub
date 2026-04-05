package interceptors

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/policies"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Auth(v *validator.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Check if method has an auth policy
		policy, exists := policies.AuthPolicies[info.FullMethod]
		if !exists {
			return handler(ctx, req)
		}

		// Check if metadata exists
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ce.NewError(ce.CodeMissingMetadata, ce.MsgInternalServer, nil)
		}

		// Check if policy requires subdomain authorization
		var subdomain string
		if policy.RequireSubdomain() {
			values := md.Get(constants.MDKeySubdomain)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeUnauthorizedSubdomain,
					ce.MsgUnauthorized,
					errors.New("subdomain missing from metadata"),
				)
			}

			subdomain = values[0]
			if !policy.IsSubdomainAuthorized(subdomain) {
				return nil, ce.NewError(
					ce.CodeUnauthorizedSubdomain,
					ce.MsgUnauthorized,
					errors.New("subdomain unauthorized"),
					logger.NewField("subdomain", subdomain),
				)
			}

			ctx = context.WithValue(ctx, constants.CtxKeySubdomain, subdomain)
		}

		// Check if policy requires authentication
		var authID int64
		if policy.RequireAuth() {
			values := md.Get(constants.MDKeyAuthID)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeUnauthenticated,
					ce.MsgUnauthenticated,
					errors.New("auth_id missing from metadata"),
				)
			}

			id, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return nil, ce.NewError(
					ce.CodeTypeConversionFailed,
					ce.MsgInternalServer,
					fmt.Errorf("failed to convert auth_id (%v) to int64: %w", values[0], err),
				)
			}

			authID = id
		}

		// Check if policy requires school information
		var schoolID int64
		if policy.RequireSchool() {
			values := md.Get(constants.MDKeySchoolID)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeMissingMetadata,
					ce.MsgInternalServer,
					errors.New("school_id missing from metadata"),
				)
			}

			id, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return nil, ce.NewError(
					ce.CodeTypeConversionFailed,
					ce.MsgInternalServer,
					fmt.Errorf("failed to convert school_id (%v) to int64: %w", values[0], err),
				)
			}

			schoolID = id
		}

		// Check if policy requires role authorization
		var role string
		if policy.RequireRole() {
			values := md.Get(constants.MDKeyRole)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeUnauthorizedRole,
					ce.MsgUnauthorized,
					errors.New("role missing from metadata"),
				)
			}

			role = values[0]
			if !policy.IsRoleAuthorized(role) {
				return nil, ce.NewError(
					ce.CodeUnauthorizedRole,
					ce.MsgUnauthorized,
					errors.New("role unauthorized"),
					logger.NewField("role", role),
				)
			}
			if subdomain != "" && !v.RoleAllowedSubdomain(role, subdomain) {
				return nil, ce.NewError(
					ce.CodeUnauthorizedSubdomain,
					ce.MsgUnauthorized,
					errors.New("role unauthorized by subdomain"),
					logger.NewField("role", role),
					logger.NewField("subdomain", subdomain),
				)
			}
		}

		// Check if policy requires verification
		var isVerified bool
		if policy.RequireVerification() {
			values := md.Get(constants.MDKeyIsVerified)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeAuthNotVerified,
					ce.MsgAuthNotVerified,
					errors.New("is_verified missing from metadata"),
				)
			}

			verified, err := strconv.ParseBool(values[0])
			if err != nil {
				return nil, ce.NewError(
					ce.CodeTypeConversionFailed,
					ce.MsgInternalServer,
					fmt.Errorf("failed to convert is_verified (%v) to bool: %w", values[0], err),
				)
			}
			if !verified {
				return nil, ce.NewError(ce.CodeAuthNotVerified, ce.MsgAuthNotVerified, nil)
			}

			isVerified = verified
		}

		return handler(
			context.WithValue(
				ctx,
				constants.CtxKeyAuth,
				&models.AuthContext{
					AuthID:     authID,
					SchoolID:   schoolID,
					Role:       role,
					IsVerified: isVerified,
				},
			),
			req,
		)
	}
}
