package testcontainer

import (
	"context"
	"github.com/redis/go-redis/v9"
	testContainerRedis "github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainer struct {
	Address     string
	container   *testContainerRedis.RedisContainer
	RedisClient *redis.Client
}

func NewRedisContainer(ctx context.Context) (*RedisContainer, error) {
	redisContainer, err := testContainerRedis.Run(ctx, "docker.io/redis:7")
	if err != nil {
		return nil, err
	}

	endpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	redisUrl, err := redis.ParseURL("redis://" + endpoint)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		Address:     "redis://" + endpoint,
		container:   nil,
		RedisClient: redis.NewClient(redisUrl),
	}, nil
}

func (c RedisContainer) Shutdown(ctx context.Context) {
	_ = c.RedisClient.Close()
	_ = c.container.Terminate(ctx)
}
