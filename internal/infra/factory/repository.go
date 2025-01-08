package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/mongo"
)

func repositoryFactory(mongoDatabase *mongo.Database) *factory.RepositoryFactory {
	return factory.NewRepositoryFactory(
		repository.NewCustomerMongoRepository(mongoDatabase),
		repository.NewProductMongoRepository(mongoDatabase),
		repository.NewShoppingCartMongoRepository(mongoDatabase),
		repository.NewOrderMongoRepository(mongoDatabase),
		repository.NewAdminMongoRepository(mongoDatabase),
		repository.NewMerchantMongoRepository(mongoDatabase),
		repository.NewTransactionMongoRepository(mongoDatabase),
		repository.NewOrderStatusHistoryMongoRepository(mongoDatabase),
		repository.NewCustomerNotificationMongoRepository(mongoDatabase),
		repository.NewNotificationMongoRepository(mongoDatabase),
		repository.NewCategoriesMongoRepository(mongoDatabase),
	)
}
