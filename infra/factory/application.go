package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewApplicationFactory(database *sql.DB) *factory.ApplicationFactory {
	return &factory.ApplicationFactory{
		GatewayFactory:    gatewayFactory(),
		RepositoryFactory: repositoryFactory(database),
		TokenFactory:      tokenFactory(),
		EventFactory:      eventFactory(),
	}
}

func eventFactory() *factory.EventFactory {
	kafkaServer := os.Getenv("KAFKA_SERVER")
	eventEmitter := event.NewKafkaEventEmitter(kafkaServer)
	eventHandler := event.NewKafkaEventHandler(kafkaServer)
	return factory.NewEventFactory(eventEmitter, eventHandler)
}

func tokenFactory() *factory.TokenFactory {
	privateKey := loadPrivateKey(os.Getenv("JWE_PRIVATE_KEY_PATH"))
	tokenManager := token.NewCustomerTokenManager(privateKey)
	return factory.NewTokenFactory(tokenManager)
}

func loadPrivateKey(file string) *rsa.PrivateKey {
	pemData, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("failed to decode PEM block containing private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

func gatewayFactory() *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(os.Getenv("GOOGLE_AUTH_URL")),
	)
}

func repositoryFactory(database *sql.DB) *factory.RepositoryFactory {
	connection := repository.NewDefaultSqlConnection(database)
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(options)
	customerRepository := repository.NewDefaultCustomerRepository(connection)
	productRepository := repository.NewDefaultProductRepository(connection)
	shoppingCartRepository := repository.NewShoppingCartRedisRepository(redisClient)
	return factory.NewRepositoryFactory(customerRepository, productRepository, shoppingCartRepository)
}
