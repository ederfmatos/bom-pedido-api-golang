package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/lock"
	"bom-pedido-api/pkg/mongo"
	"github.com/redis/go-redis/v9"
)

func NewApplicationFactory(environment *config.Environment, redisClient *redis.Client, mongoDatabase *mongo.Database) *factory.ApplicationFactory {
	locker := lock.NewRedisLocker(redisClient)
	return factory.NewApplicationFactory(
		gatewayFactory(environment, mongoDatabase),
		repositoryFactory(mongoDatabase),
		tokenFactory(environment),
		eventFactory(environment, locker, mongoDatabase),
		queryFactory(mongoDatabase),
		locker,
	)
}
