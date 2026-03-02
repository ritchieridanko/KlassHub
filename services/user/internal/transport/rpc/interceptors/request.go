package interceptors

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/constants"
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

		var requestID string
		if values := md.Get(constants.MDKeyRequestID); len(values) > 0 {
			requestID = values[0]
		}
		return handler(
			context.WithValue(
				ctx,
				constants.CtxKeyRequestID,
				requestID,
			),
			req,
		)
	}
}
