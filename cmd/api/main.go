package main

import (
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/http"
	"bom-pedido-api/internal/infra/messaging"
	"fmt"
	"log"
)

func main() {
	config.ConfigureLogger()
	environment, err := config.LoadEnvironment()
	if err != nil {
		log.Fatalf("load environment: %v", err)
	}

	redisClient := config.Redis(environment.RedisUrl)
	mongoClient := config.Mongo(environment.MongoUrl)

	applicationFactory := factory.NewApplicationFactory(environment, redisClient, mongoClient)
	defer applicationFactory.Close()

	go messaging.HandleEvents(applicationFactory)

	server := http.NewServer(redisClient, mongoClient, environment)
	server.ConfigureRoutes(applicationFactory)
	go server.Run(fmt.Sprintf(":%s", environment.Port))
	server.AwaitInterruptSignal()
	server.Shutdown()
}
