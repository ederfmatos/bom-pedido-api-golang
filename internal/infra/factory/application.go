package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/lock"
	mongo2 "bom-pedido-api/pkg/mongo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApplicationFactory(environment *config.Environment, redisClient *redis.Client, mongoClient *mongo.Client) *factory.ApplicationFactory {
	locker := lock.NewRedisLocker(redisClient)
	mongoDatabase := mongo2.NewDatabase(mongoClient.Database(environment.MongoDatabaseName))
	return factory.NewApplicationFactory(
		gatewayFactory(environment, mongoDatabase),
		repositoryFactory(mongoDatabase),
		tokenFactory(environment),
		eventFactory(environment, locker, mongoDatabase),
		queryFactory(mongoDatabase),
		locker,
	)
}
