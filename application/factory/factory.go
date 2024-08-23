package factory

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
)

type ApplicationFactory struct {
	*GatewayFactory
	*RepositoryFactory
	*TokenFactory
	*EventFactory
	*QueryFactory
	Locker       lock.Locker
	EmailGateway gateway.EmailGateway
}

func (factory *ApplicationFactory) Close() {
	factory.EventHandler.Close()
}

func NewApplicationFactory(
	gatewayFactory *GatewayFactory,
	repositoryFactory *RepositoryFactory,
	tokenFactory *TokenFactory,
	eventFactory *EventFactory,
	queryFactory *QueryFactory,
	locker lock.Locker,
	emailGateway gateway.EmailGateway,
) *ApplicationFactory {
	return &ApplicationFactory{
		GatewayFactory:    gatewayFactory,
		RepositoryFactory: repositoryFactory,
		TokenFactory:      tokenFactory,
		EventFactory:      eventFactory,
		QueryFactory:      queryFactory,
		Locker:            locker,
		EmailGateway:      emailGateway,
	}
}
