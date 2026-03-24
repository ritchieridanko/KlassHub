package interceptors

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
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

		var ua, ip string
		if values := md.Get(constants.MDKeyUserAgent); len(values) > 0 {
			ua = values[0]
		}
		if values := md.Get(constants.MDKeyIPAddress); len(values) > 0 {
			ip = values[0]
		}

		return handler(
			context.WithValue(
				ctx,
				constants.CtxKeyTransport,
				&models.TransportContext{
					UserAgent: ua,
					IPAddress: ip,
				},
			),
			req,
		)
	}
}
