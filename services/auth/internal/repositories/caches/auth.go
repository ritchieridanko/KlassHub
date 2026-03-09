package caches

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type AuthCache interface {
	IsEmailReserved(ctx context.Context, email string) (exists bool, err *ce.Error)
}

type authCache struct {
	cache *cache.Cache
}

func NewAuthCache(cc *cache.Cache) AuthCache {
	return &authCache{cache: cc}
}

func (c *authCache) IsEmailReserved(ctx context.Context, email string) (bool, *ce.Error) {
	exists, err := c.cache.Exists(
		ctx,
		fmt.Sprintf(
			"%s:%s",
			constants.CachePrefixEmailReservation,
			email,
		),
	)
	if err != nil {
		return false, ce.NewError(
			ce.CodeCacheCommandExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to check if email is reserved: %w", err),
		)
	}
	return exists, nil
}
