package cache

import (
	"log/slog"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	*cache.Cache
	logger *slog.Logger
}

func NewCacheService(logger *slog.Logger, redisClient *redis.Client) (*CacheService, error) {
	cache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return &CacheService{
		Cache:  cache,
		logger: logger.With("service", "cache"),
	}, nil
}
