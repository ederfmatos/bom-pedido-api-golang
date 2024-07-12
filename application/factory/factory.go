package factory

import "bom-pedido-api/application/lock"

type ApplicationFactory struct {
	*GatewayFactory
	*RepositoryFactory
	*TokenFactory
	*EventFactory
	*QueryFactory
	Locker lock.Locker
}

func NewApplicationFactory(
	gatewayFactory *GatewayFactory,
	repositoryFactory *RepositoryFactory,
	tokenFactory *TokenFactory,
	eventFactory *EventFactory,
	queryFactory *QueryFactory,
	locker lock.Locker,
) *ApplicationFactory {
	return &ApplicationFactory{
		GatewayFactory:    gatewayFactory,
		RepositoryFactory: repositoryFactory,
		TokenFactory:      tokenFactory,
		EventFactory:      eventFactory,
		QueryFactory:      queryFactory,
		Locker:            locker,
	}
}
