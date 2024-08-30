package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway/email"
	"bom-pedido-api/infra/gateway/google"
	"bom-pedido-api/infra/gateway/pix"
	"bom-pedido-api/infra/lock"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/test"
	"bom-pedido-api/infra/token"
)

func NewTestApplicationFactory() *factory.ApplicationFactory {
	return factory.NewApplicationFactory(
		factory.NewGatewayFactory(google.NewFakeGoogleGateway(), pix.NewFakePixGateway()),
		factory.NewRepositoryFactory(
			repository.NewCustomerMemoryRepository(),
			repository.NewProductMemoryRepository(),
			repository.NewShoppingCartMemoryRepository(),
			repository.NewOrderMemoryRepository(),
			repository.NewAdminMemoryRepository(),
			repository.NewMerchantMemoryRepository(),
			repository.NewTransactionMemoryRepository(),
		),
		factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		factory.NewEventFactory(event.NewMemoryEventHandler()),
		nil,
		lock.NewMemoryLocker(),
		email.NewFakeEmailGateway(),
	)
}

func NewContainerApplicationFactory(container *test.Container) *factory.ApplicationFactory {
	sqlConnection := repository.NewDefaultSqlConnection(container.Database)
	return factory.NewApplicationFactory(
		factory.NewGatewayFactory(google.NewFakeGoogleGateway(), pix.NewFakePixGateway()),
		factory.NewRepositoryFactory(
			repository.NewDefaultCustomerRepository(sqlConnection),
			repository.NewDefaultProductRepository(sqlConnection),
			repository.NewShoppingCartMongoRepository(container.MongoDatabase),
			repository.NewDefaultOrderRepository(sqlConnection),
			repository.NewDefaultAdminRepository(sqlConnection),
			repository.NewDefaultMerchantRepository(sqlConnection),
			repository.NewDefaultTransactionRepository(sqlConnection),
		),
		factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		factory.NewEventFactory(event.NewMemoryEventHandler()),
		queryFactory(sqlConnection),
		lock.NewRedisLocker(container.RedisClient),
		email.NewFakeEmailGateway(),
	)
}
