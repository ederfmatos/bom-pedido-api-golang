package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func repositoryFactory(connection repository.SqlConnection, mongoDatabase *mongo.Database) *factory.RepositoryFactory {
	return factory.NewRepositoryFactory(
		repository.NewDefaultCustomerRepository(connection),
		repository.NewDefaultProductRepository(connection),
		repository.NewShoppingCartMongoRepository(mongoDatabase),
		repository.NewDefaultOrderRepository(connection),
		repository.NewDefaultAdminRepository(connection),
		repository.NewDefaultMerchantRepository(connection),
		repository.NewDefaultTransactionRepository(connection),
		repository.NewDefaultOrderStatusHistoryRepository(connection),
		repository.NewCustomerNotificationMongoRepository(mongoDatabase),
		repository.NewNotificationMongoRepository(mongoDatabase),
		repository.NewDefaultProductCategoryRepository(connection),
	)
}
