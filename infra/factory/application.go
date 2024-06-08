package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"database/sql"
	"os"
)

func NewApplicationFactory(database *sql.DB) *factory.ApplicationFactory {
	return &factory.ApplicationFactory{
		GatewayFactory:    gatewayFactory(),
		RepositoryFactory: repositoryFactory(database),
		TokenFactory:      tokenFactory(),
		EventFactory:      eventFactory(),
	}
}

func eventFactory() *factory.EventFactory {
	kafkaServer := os.Getenv("KAFKA_SERVER")
	eventEmitter := event.NewKafkaEventEmitter(kafkaServer)
	eventHandler := event.NewKafkaEventHandler(kafkaServer)
	return factory.NewEventFactory(eventEmitter, eventHandler)
}

func tokenFactory() *factory.TokenFactory {
	return factory.NewTokenFactory(nil)
}

func gatewayFactory() *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(os.Getenv("GOOGLE_AUTH_URL")),
	)
}

func repositoryFactory(database *sql.DB) *factory.RepositoryFactory {
	connection := repository.NewDefaultSqlConnection(database)
	customerRepository := repository.NewDefaultCustomerRepository(connection)
	productRepository := repository.NewDefaultProductRepository(connection)
	return factory.NewRepositoryFactory(customerRepository, productRepository)
}
