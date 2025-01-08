package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/config"
	infraEvent "bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/mongo"
)

func eventFactory(environment *config.Environment, locker lock.Locker, mongoDatabase *mongo.Database) *factory.EventFactory {
	eventHandler := infraEvent.NewRabbitMqAdapter(environment.RabbitMqServer)
	outboxCollection := mongoDatabase.ForCollection("outbox")
	outboxRepository := repository.NewMongoOutboxRepository(outboxCollection)
	handler := infraEvent.NewOutboxEventHandler(eventHandler, outboxRepository, outboxCollection, locker)
	return factory.NewEventFactory(handler)
}
