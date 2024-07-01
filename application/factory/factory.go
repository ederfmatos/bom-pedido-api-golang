package factory

type ApplicationFactory struct {
	*GatewayFactory
	*RepositoryFactory
	*TokenFactory
	*EventFactory
	*QueryFactory
}

func NewApplicationFactory(gatewayFactory *GatewayFactory, repositoryFactory *RepositoryFactory, tokenFactory *TokenFactory, eventFactory *EventFactory, queryFactory *QueryFactory) *ApplicationFactory {
	return &ApplicationFactory{
		GatewayFactory:    gatewayFactory,
		RepositoryFactory: repositoryFactory,
		TokenFactory:      tokenFactory,
		EventFactory:      eventFactory,
		QueryFactory:      queryFactory,
	}
}
