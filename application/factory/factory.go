package factory

import (
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
)

type ApplicationFactory struct {
	*GatewayFactory
	*RepositoryFactory
	*TokenFactory
}

func NewTestApplicationFactory() *ApplicationFactory {
	return &ApplicationFactory{
		&GatewayFactory{GoogleGateway: gateway.NewFakeGoogleGateway()},
		&RepositoryFactory{CustomerRepository: repository.NewCustomerMemoryRepository()},
		&TokenFactory{CustomerTokenManager: token.NewFakeCustomerTokenManager()},
	}
}
