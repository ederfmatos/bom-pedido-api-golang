package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/test"
	"bom-pedido-api/infra/token"
)

func NewTestApplicationFactory() *factory.ApplicationFactory {
	return factory.NewApplicationFactory(
		factory.NewGatewayFactory(gateway.NewFakeGoogleGateway()),
		factory.NewRepositoryFactory(
			repository.NewCustomerMemoryRepository(),
			repository.NewProductMemoryRepository(),
			repository.NewShoppingCartMemoryRepository(),
			repository.NewOrderMemoryRepository(),
		),
		factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		factory.NewEventFactory(event.NewMemoryEventHandler()),
		nil, nil,
	)
}

func NewContainerApplicationFactory(container *test.Container) *factory.ApplicationFactory {
	sqlConnection := repository.NewDefaultSqlConnection(container.Database)
	mongoDatabase := container.MongoDatabase
	return factory.NewApplicationFactory(
		factory.NewGatewayFactory(gateway.NewFakeGoogleGateway()),
		factory.NewRepositoryFactory(
			repository.NewDefaultCustomerRepository(sqlConnection),
			repository.NewDefaultProductRepository(sqlConnection),
			repository.NewShoppingCartMongoRepository(mongoDatabase),
			repository.NewDefaultOrderRepository(sqlConnection),
		),
		factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		factory.NewEventFactory(event.NewMemoryEventHandler()),
		nil, nil,
	)
}
