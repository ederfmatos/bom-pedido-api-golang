package main

import (
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/presentation/http"
	"bom-pedido-api/presentation/messaging"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
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

	applicationFactory := factory.NewApplicationFactory(database, environment, redisClient)
	defer applicationFactory.EventHandler.Close()

	go messaging.HandleEvents(applicationFactory)

	server := http.Server(applicationFactory)
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%s", environment.Port)))
}
