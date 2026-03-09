package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{client: c}
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	res, err := c.client.Exists(ctx, key).Result()
	return res > 0, err
}
