package factory

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/config"
	infraEvent "bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/repository/outbox"
	"go.mongodb.org/mongo-driver/mongo"
)

func eventFactory(environment *config.Environment, locker lock.Locker, mongoDatabase *mongo.Database) *factory.EventFactory {
	eventHandler := makeEventHandler(environment)
	collection := mongoDatabase.Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := infraEvent.NewMongoStream(collection)
	handler := infraEvent.NewOutboxEventHandler(eventHandler, outboxRepository, mongoStream, locker)
	return factory.NewEventFactory(handler)
}

func makeEventHandler(environment *config.Environment) event.Handler {
	switch environment.MessagingStrategy {
	case "KAFKA":
		return infraEvent.NewKafkaEventHandler(environment)
	default:
		return infraEvent.NewRabbitMqAdapter(environment.RabbitMqServer)
	}
}
