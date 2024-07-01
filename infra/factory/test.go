package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
)

func NewTestApplicationFactory() *factory.ApplicationFactory {
	return &factory.ApplicationFactory{
		GatewayFactory: factory.NewGatewayFactory(gateway.NewFakeGoogleGateway()),
		RepositoryFactory: factory.NewRepositoryFactory(
			repository.NewCustomerMemoryRepository(),
			repository.NewProductMemoryRepository(),
			repository.NewShoppingCartMemoryRepository(),
		),
		TokenFactory: factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		EventFactory: factory.NewEventFactory(event.NewMemoryEventEmitter(), nil),
	}
}
