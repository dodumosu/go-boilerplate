package redisutil

import (
	"go-boilerplate/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	connectionString, err := cfg.BuildConnectionString()
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	return client, nil
}
