package factory

import (
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
)

type ApplicationFactory struct {
	*GatewayFactory
	*RepositoryFactory
	*TokenFactory
	*EventFactory
}

func NewTestApplicationFactory() *ApplicationFactory {
	return &ApplicationFactory{
		NewGatewayFactory(gateway.NewFakeGoogleGateway()),
		NewRepositoryFactory(
			repository.NewCustomerMemoryRepository(),
			repository.NewProductMemoryRepository(),
		),
		NewTokenFactory(token.NewFakeCustomerTokenManager()),
		NewEventFactory(event.NewMemoryEventEmitter(), nil),
	}
}
