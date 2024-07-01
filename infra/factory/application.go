package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/query"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewApplicationFactory(database *sql.DB, environment *env.Environment) *factory.ApplicationFactory {
	connection := repository.NewDefaultSqlConnection(database)
	return &factory.ApplicationFactory{
		GatewayFactory:    gatewayFactory(environment),
		RepositoryFactory: repositoryFactory(connection, environment),
		TokenFactory:      tokenFactory(environment),
		EventFactory:      eventFactory(environment),
		QueryFactory:      queryFactory(connection),
	}
}

func queryFactory(connection repository.SqlConnection) *factory.QueryFactory {
	return factory.NewQueryFactory(query.NewProductSqlQuery(connection))
}

func eventFactory(environment *env.Environment) *factory.EventFactory {
	rabbitMqAdapter := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	return factory.NewEventFactory(rabbitMqAdapter, rabbitMqAdapter)
}

func tokenFactory(environment *env.Environment) *factory.TokenFactory {
	privateKey := loadPrivateKey(environment.JwePrivateKeyPath)
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

func gatewayFactory(environment *env.Environment) *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(environment.GoogleAuthUrl),
	)
}

func repositoryFactory(connection repository.SqlConnection, environment *env.Environment) *factory.RepositoryFactory {
	options, err := redis.ParseURL(environment.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(options)
	customerRepository := repository.NewDefaultCustomerRepository(connection)
	productRepository := repository.NewDefaultProductRepository(connection)
	shoppingCartRepository := repository.NewShoppingCartRedisRepository(redisClient)
	return factory.NewRepositoryFactory(customerRepository, productRepository, shoppingCartRepository)
}
