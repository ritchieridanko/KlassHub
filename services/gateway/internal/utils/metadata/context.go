package metadata

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"google.golang.org/grpc/metadata"
)

func ToOutgoingCtx(ctx context.Context, pairs ...Pair) context.Context {
	defaultSize := 2
	kv := append(
		make([]string, 0, (len(pairs)*2)+defaultSize),
		constants.MDKeyRequestID,
		utils.CtxRequestID(ctx),
	)

	for _, pair := range pairs {
		kv = append(kv, pair.key, pair.value)
	}
	return metadata.AppendToOutgoingContext(ctx, kv...)
}
