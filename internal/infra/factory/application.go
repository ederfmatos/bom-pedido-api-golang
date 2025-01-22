package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/lock"
	"bom-pedido-api/pkg/mongo"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewApplicationFactory(environment *config.Environment, redisClient *redis.Client, mongoDatabase *mongo.Database) (*factory.ApplicationFactory, error) {
	gateway, err := gatewayFactory(environment, mongoDatabase)
	if err != nil {
		return nil, fmt.Errorf("create gateway factory: %v", err)
	}

	token, err := tokenFactory(environment)
	if err != nil {
		return nil, fmt.Errorf("create token factory: %v", err)
	}

	locker := lock.NewRedisLocker(redisClient)
	event, err := eventFactory(environment, locker, mongoDatabase)
	if err != nil {
		return nil, fmt.Errorf("create event factory: %v", err)
	}

	repository := repositoryFactory(mongoDatabase)
	query := queryFactory(mongoDatabase)
	return factory.NewApplicationFactory(gateway, repository, token, event, query, locker), nil
}
