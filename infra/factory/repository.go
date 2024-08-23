package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func repositoryFactory(connection repository.SqlConnection, mongoDatabase *mongo.Database) *factory.RepositoryFactory {
	customerRepository := repository.NewDefaultCustomerRepository(connection)
	productRepository := repository.NewDefaultProductRepository(connection)
	orderRepository := repository.NewDefaultOrderRepository(connection)
	shoppingCartRepository := repository.NewShoppingCartMongoRepository(mongoDatabase)
	adminRepository := repository.NewDefaultAdminRepository(connection)
	return factory.NewRepositoryFactory(customerRepository, productRepository, shoppingCartRepository, orderRepository, adminRepository)
}
