package caches

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type TokenCache interface {
	CreateVerification(ctx context.Context, data *models.CreateVerificationTokenData) (err *ce.Error)
	UseVerification(ctx context.Context, token string) (authID int64, err *ce.Error)
}

type tokenCache struct {
	config *configs.Auth
	cache  *cache.Cache
}

func NewTokenCache(cfg *configs.Auth, cc *cache.Cache) TokenCache {
	return &tokenCache{config: cfg, cache: cc}
}

func (c *tokenCache) CreateVerification(ctx context.Context, data *models.CreateVerificationTokenData) *ce.Error {
	prefix := constants.CachePrefixEmailVerification
	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end

		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "s:crever", script,
		[]string{
			fmt.Sprintf("%s:%d", prefix, data.AuthID),
			fmt.Sprintf("%s:%s", prefix, data.Token),
			prefix,
		},
		data.Token, data.AuthID, int(data.Duration.Seconds()),
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create verification token: %w", err),
			logger.NewField("auth_id", data.AuthID),
		)
	}

	return nil
}

func (c *tokenCache) UseVerification(ctx context.Context, token string) (int64, *ce.Error) {
	prefix := constants.CachePrefixEmailVerification
	script := `
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`

	res, err := c.cache.Evaluate(
		ctx, "s:usver", script,
		[]string{
			fmt.Sprintf("%s:%s", prefix, token),
			prefix,
		},
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to use verification token: %w", err)

		if errors.Is(err, ce.ErrCacheNoResult) {
			return 0, ce.NewError(ce.CodeTokenNotFound, ce.MsgInvalidToken, wrappedErr)
		}
		return 0, ce.NewError(ce.CodeCacheScriptExec, ce.MsgInternalServer, wrappedErr)
	}

	authID, err := utils.ToInt64(res)
	if err != nil {
		return 0, ce.NewError(
			ce.CodeTypeConversionFailed,
			ce.MsgInternalServer,
			fmt.Errorf("failed to use verification token: %w", err),
		)
	}

	return authID, nil
}
