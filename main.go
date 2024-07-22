package main

import (
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/http"
	"bom-pedido-api/infra/messaging"
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
	redisUrl, err := redis.ParseURL(environment.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisUrl)

	clientOptions := options.Client().ApplyURI(environment.MongoUrl)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	applicationFactory := factory.NewApplicationFactory(database, environment, redisClient, mongoClient)
	defer applicationFactory.EventHandler.Close()

	go messaging.HandleEvents(applicationFactory)

	server := http.NewServer(database, redisClient, mongoClient, environment)
	server.ConfigureRoutes(applicationFactory)
	go server.Run(fmt.Sprintf(":%s", environment.Port))
	server.AwaitInterruptSignal()
	server.Shutdown()
}
