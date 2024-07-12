package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/lock"
	"bom-pedido-api/infra/repository/outbox"
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func eventFactory(environment *env.Environment, redisClient *redis.Client) *factory.EventFactory {
	rabbitMqAdapter := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	clientOptions := options.Client().ApplyURI(environment.MongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	locker := lock.NewRedisLocker(redisClient)
	collection := client.Database(environment.MongoDatabaseName).Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := event.NewMongoStream(collection)
	handler := event.NewOutboxEventHandler(rabbitMqAdapter, outboxRepository, mongoStream, locker)
	return factory.NewEventFactory(handler)
}
