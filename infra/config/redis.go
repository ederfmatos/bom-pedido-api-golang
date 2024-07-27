package config

import "github.com/redis/go-redis/v9"

func Redis(url string) *redis.Client {
	redisUrl, err := redis.ParseURL(url)
	failOnError(err, "Failed to parse redis url")
	return redis.NewClient(redisUrl)
}
