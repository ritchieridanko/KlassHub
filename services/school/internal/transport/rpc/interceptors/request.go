package interceptors

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		var requestID, userAgent, ipAddress string
		if values := md.Get(constants.MDKeyRequestID); len(values) > 0 {
			requestID = values[0]
		}
		if values := md.Get(constants.MDKeyUserAgent); len(values) > 0 {
			userAgent = values[0]
		}
		if values := md.Get(constants.MDKeyIPAddress); len(values) > 0 {
			ipAddress = values[0]
		}

		ctx = context.WithValue(ctx, constants.CtxKeyRequestID, requestID)
		ctx = context.WithValue(ctx, constants.CtxKeyUserAgent, userAgent)
		ctx = context.WithValue(ctx, constants.CtxKeyIPAddress, ipAddress)

		return handler(ctx, req)
	}
}
