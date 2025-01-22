package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/mongo"
	"fmt"
)

func eventFactory(environment *config.Environment, locker lock.Locker, mongoDatabase *mongo.Database) (*factory.EventFactory, error) {
	eventHandler, err := event.NewRabbitMqEventHandler(environment.RabbitMqServer)
	if err != nil {
		return nil, fmt.Errorf("create rabbitmq event handler: %v", err)
	}

	outboxRepository := repository.NewMongoOutboxRepository(mongoDatabase)
	handler, err := event.NewOutboxEventHandler(eventHandler, outboxRepository, locker)
	if err != nil {
		return nil, fmt.Errorf("create outbox event handler: %v", err)
	}

	return factory.NewEventFactory(handler), nil
}
