package main

import (
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/presentation/http"
	"bom-pedido-api/presentation/http/health"
	"bom-pedido-api/presentation/messaging"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	environment := env.LoadEnvironment(".env")
	database, err := sql.Open(environment.DatabaseDriver, environment.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	redisUrl, err := redis.ParseURL(environment.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisUrl)
	defer redisClient.Close()

	clientOptions := options.Client().ApplyURI(environment.MongoUrl)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	applicationFactory := factory.NewApplicationFactory(database, environment, redisClient, mongoClient)
	defer applicationFactory.EventHandler.Close()

	go messaging.HandleEvents(applicationFactory)

	server := http.Server(applicationFactory)
	server.GET("/api/health", health.Handle(database, redisClient, mongoClient))

	server.Logger.Fatal(server.Start(fmt.Sprintf(":%s", environment.Port)))
}
