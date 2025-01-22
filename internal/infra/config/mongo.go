package config

import (
	mongo2 "bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func Mongo(url, database string) (*mongo2.Database, error) {
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
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %v", err)
	}

	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %v", err)
	}

	return mongo2.NewDatabase(mongoClient.Database(database)), nil
}
