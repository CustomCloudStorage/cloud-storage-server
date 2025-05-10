package databases

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	RPort     string
	RPassword string
}

func GetRedis(cfg Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RPort,
		Password: cfg.RPassword,
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}
