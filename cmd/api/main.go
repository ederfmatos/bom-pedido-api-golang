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
	server, err := makeServer()
	if err != nil {
		log.Fatalf("make server: %v", err)
	}

	go server.Run()
	server.AwaitInterruptSignal()
	server.Shutdown()
}

func makeServer() (*http.Server, error) {
	config.ConfigureLogger()
	environment, err := config.LoadEnvironment()
	if err != nil {
		return nil, fmt.Errorf("load environment: %v", err)
	}

	redisClient, err := config.Redis(environment.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("connect redis: %v", err)
	}

	mongoDatabase, err := config.Mongo(environment.MongoUrl, environment.MongoDatabaseName)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %v", err)
	}

	applicationFactory, err := factory.NewApplicationFactory(environment, redisClient, mongoDatabase)
	if err != nil {
		return nil, fmt.Errorf("create application factory: %v", err)
	}

	go messaging.HandleEvents(applicationFactory)

	server := http.NewServer(redisClient, mongoDatabase, environment, applicationFactory)
	server.ConfigureRoutes()

	return server, nil
}
