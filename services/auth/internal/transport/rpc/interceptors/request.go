package interceptors

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Request() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		if values := md.Get(constants.MDKeyRequestID); len(values) > 0 {
			ctx = context.WithValue(ctx, constants.CtxKeyRequestID, values[0])
		}
		if values := md.Get(constants.MDKeyUserAgent); len(values) > 0 {
			ctx = context.WithValue(ctx, constants.CtxKeyUserAgent, values[0])
		}
		if values := md.Get(constants.MDKeyIPAddress); len(values) > 0 {
			ctx = context.WithValue(ctx, constants.CtxKeyIPAddress, values[0])
		}
		if values := md.Get(constants.MDKeySubdomain); len(values) > 0 {
			ctx = context.WithValue(ctx, constants.CtxKeySubdomain, values[0])
		}
		return handler(ctx, req)
	}
}
