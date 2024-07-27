package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Mongo(url string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(url)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	failOnError(err, "Failed to connect to mongo")
	err = mongoClient.Ping(context.TODO(), nil)
	failOnError(err, "Failed to ping mongo")
	return mongoClient
}
