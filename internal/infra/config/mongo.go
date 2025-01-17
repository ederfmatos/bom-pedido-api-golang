package config

import (
	mongo2 "bom-pedido-api/pkg/mongo"
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"log/slog"
)

func Mongo(url, database string) *mongo2.Database {
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
	return mongo2.NewDatabase(mongoClient.Database(database))
}
