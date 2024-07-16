package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/repository/outbox"
	"go.mongodb.org/mongo-driver/mongo"
)

func eventFactory(environment *env.Environment, locker lock.Locker, mongoClient *mongo.Client) *factory.EventFactory {
	//eventHandler := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	eventHandler := event.NewKafkaEventHandler(environment)
	collection := mongoClient.Database(environment.MongoDatabaseName).Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := event.NewMongoStream(collection)
	handler := event.NewOutboxEventHandler(eventHandler, outboxRepository, mongoStream, locker)
	return factory.NewEventFactory(handler)
}
