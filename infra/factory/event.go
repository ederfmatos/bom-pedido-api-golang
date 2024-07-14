package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/repository/outbox"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func eventFactory(environment *env.Environment, locker lock.Locker) *factory.EventFactory {
	//eventHandler := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	eventHandler := event.NewKafkaEventHandler(environment)
	clientOptions := options.Client().ApplyURI(environment.MongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	collection := client.Database(environment.MongoDatabaseName).Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := event.NewMongoStream(collection)
	handler := event.NewOutboxEventHandler(eventHandler, outboxRepository, mongoStream, locker)
	return factory.NewEventFactory(handler)
}
