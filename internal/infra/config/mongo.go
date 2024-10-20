package config

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"log/slog"
)

func Mongo(url string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(url)
	monitor := otelmongo.NewMonitor()
	started := monitor.Started
	monitor.Started = func(ctx context.Context, event *event.CommandStartedEvent) {
		if event.CommandName == "getMore" {
			return
		}
		started(ctx, event)
	}
	clientOptions.Monitor = monitor
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	failOnError(err, "Failed to connect to mongo")
	err = mongoClient.Ping(context.Background(), nil)
	failOnError(err, "Failed to ping mongo")
	slog.Info("Connected to mongo successfully")
	return mongoClient
}
