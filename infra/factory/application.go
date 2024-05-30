package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/environment"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"database/sql"
)

func NewApplicationFactory(database *sql.DB) *factory.ApplicationFactory {
	return &factory.ApplicationFactory{
		GatewayFactory:    gatewayFactory(),
		RepositoryFactory: repositoryFactory(database),
		TokenFactory:      tokenFactory(),
	}
}

func tokenFactory() *factory.TokenFactory {
	return factory.NewTokenFactory(nil)
}

func gatewayFactory() *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(environment.GetEnvVar("GOOGLE_AUTH_URL")),
	)
}

func repositoryFactory(database *sql.DB) *factory.RepositoryFactory {
	connection := repository.NewDefaultSqlConnection(database)
	return factory.NewRepositoryFactory(
		repository.NewDefaultCustomerRepository(connection),
	)
}
