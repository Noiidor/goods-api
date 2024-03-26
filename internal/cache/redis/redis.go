package redis

import (
	"context"
	"fmt"
	"goods-api/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewConnection(config *config.AppConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(config.Data.Redis.URL)
	if err != nil {
		return nil, err
	}

	conn := redis.NewClient(opt)
	result, err := conn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if result != "PONG" {
		return nil, fmt.Errorf("failed to ping redis server: %v", result)
	}

	return conn, nil
}
