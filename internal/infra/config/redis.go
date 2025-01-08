package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
)

func Redis(url string) *redis.Client {
	redisUrl, err := redis.ParseURL(url)
	failOnError(err, "Failed to parse redis url")
	client := redis.NewClient(redisUrl)
	err = client.Ping(context.Background()).Err()
	failOnError(err, "Failed to ping redis")
	slog.Info("Connected to redis successfully")
	return client
}

func failOnError(err error, s string) {
	if err != nil {
		log.Fatal(s, err)
	}
}
