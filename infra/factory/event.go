package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/repository/outbox"
	"go.mongodb.org/mongo-driver/mongo"
)

func eventFactory(environment *config.Environment, locker lock.Locker, mongoDatabase *mongo.Database) *factory.EventFactory {
	//eventHandler := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	eventHandler := event.NewKafkaEventHandler(environment)
	collection := mongoDatabase.Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := event.NewMongoStream(collection)
	handler := event.NewOutboxEventHandler(eventHandler, outboxRepository, mongoStream, locker)
	return factory.NewEventFactory(handler)
}
