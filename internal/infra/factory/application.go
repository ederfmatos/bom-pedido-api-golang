package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/lock"
	"bom-pedido-api/internal/infra/repository"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApplicationFactory(database *sql.DB, environment *config.Environment, redisClient *redis.Client, mongoClient *mongo.Client) *factory.ApplicationFactory {
	connection := repository.NewDefaultSqlConnection(database)
	locker := lock.NewRedisLocker(redisClient)
	mongoDatabase := mongoClient.Database(environment.MongoDatabaseName)
	return factory.NewApplicationFactory(
		gatewayFactory(environment, connection),
		repositoryFactory(connection, mongoDatabase),
		tokenFactory(environment),
		eventFactory(environment, locker, mongoDatabase),
		queryFactory(connection),
		locker,
	)
}
