package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Redis(url string) (*redis.Client, error) {
	redisUrl, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %v", err)
	}

	client := redis.NewClient(redisUrl)
	if err = client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %v", err)
	}

	return client, nil
}
