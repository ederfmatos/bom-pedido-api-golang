package main

import (
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/http"
	"bom-pedido-api/infra/messaging"
	"fmt"
)

func main() {
	config.ConfigureLogger()
	environment := config.LoadEnvironment()
	database := config.Database(environment.DatabaseDriver, environment.DatabaseUrl)
	redisClient := config.Redis(environment.RedisUrl)
	mongoClient := config.Mongo(environment.MongoUrl)

	applicationFactory := factory.NewApplicationFactory(database, environment, redisClient, mongoClient)
	defer applicationFactory.Close()

	go messaging.HandleEvents(applicationFactory)

	server := http.NewServer(database, redisClient, mongoClient, environment)
	server.ConfigureRoutes(applicationFactory)
	go server.Run(fmt.Sprintf(":%s", environment.Port))
	server.AwaitInterruptSignal()
	server.Shutdown()
}
