package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/mongo"
)

func eventFactory(environment *config.Environment, locker lock.Locker, mongoDatabase *mongo.Database) *factory.EventFactory {
	// TODO: Handle error
	eventHandler, _ := event.NewRabbitMqEventHandler(environment.RabbitMqServer)
	outboxRepository := repository.NewMongoOutboxRepository(mongoDatabase)
	// TODO: Handle error
	handler, _ := event.NewOutboxEventHandler(eventHandler, outboxRepository, locker)
	return factory.NewEventFactory(handler)
}
