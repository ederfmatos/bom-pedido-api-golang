package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/repository"
	"github.com/redis/go-redis/v9"
)

func repositoryFactory(connection repository.SqlConnection, redisClient *redis.Client) *factory.RepositoryFactory {
	customerRepository := repository.NewDefaultCustomerRepository(connection)
	productRepository := repository.NewDefaultProductRepository(connection)
	orderRepository := repository.NewDefaultOrderRepository(connection)
	shoppingCartRepository := repository.NewShoppingCartRedisRepository(redisClient)
	return factory.NewRepositoryFactory(customerRepository, productRepository, shoppingCartRepository, orderRepository)
}
