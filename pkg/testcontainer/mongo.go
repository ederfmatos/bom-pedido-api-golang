package testcontainer

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoContainer struct {
	container   *mongodb.MongoDBContainer
	MongoClient *mongo.Client
	Address     string
}

func NewMongoContainer(ctx context.Context) (*MongoContainer, error) {
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		return nil, err
	}

	endpoint, err := mongodbContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + endpoint)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoContainer{
		MongoClient: mongoClient,
		Address:     "mongodb://" + endpoint,
		container:   mongodbContainer,
	}, nil
}

func (c MongoContainer) Shutdown(ctx context.Context) {
	if c.container != nil {
		_ = c.container.Terminate(ctx)
	}

	if c.MongoClient != nil {
		_ = c.MongoClient.Disconnect(ctx)
	}
}
