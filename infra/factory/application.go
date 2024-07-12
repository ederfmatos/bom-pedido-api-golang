package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/lock"
	"bom-pedido-api/infra/repository"
	"database/sql"
	"github.com/redis/go-redis/v9"
)

func NewApplicationFactory(database *sql.DB, environment *env.Environment, redisClient *redis.Client) *factory.ApplicationFactory {
	connection := repository.NewDefaultSqlConnection(database)
	locker := lock.NewRedisLocker(redisClient)
	return factory.NewApplicationFactory(
		gatewayFactory(environment),
		repositoryFactory(connection, redisClient),
		tokenFactory(environment),
		eventFactory(environment, locker),
		queryFactory(connection),
		locker,
	)
}
