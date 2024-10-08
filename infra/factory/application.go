package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/lock"
	"bom-pedido-api/infra/repository"
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
