package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/query"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/repository/outbox"
	"bom-pedido-api/infra/token"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func NewApplicationFactory(database *sql.DB, environment *env.Environment) *factory.ApplicationFactory {
	connection := repository.NewDefaultSqlConnection(database)
	return factory.NewApplicationFactory(
		gatewayFactory(environment),
		repositoryFactory(connection, environment),
		tokenFactory(environment),
		eventFactory(environment),
		queryFactory(connection),
	)
}

func queryFactory(connection repository.SqlConnection) *factory.QueryFactory {
	return factory.NewQueryFactory(query.NewProductSqlQuery(connection))
}

func eventFactory(environment *env.Environment) *factory.EventFactory {
	rabbitMqAdapter := event.NewRabbitMqAdapter(environment.RabbitMqServer)
	clientOptions := options.Client().ApplyURI(environment.MongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	collection := client.Database(environment.MongoDatabaseName).Collection(environment.MongoOutboxCollectionName)
	outboxRepository := outbox.NewMongoOutboxRepository(collection)
	mongoStream := event.NewMongoStream(collection)
	handler := event.NewOutboxEventHandler(rabbitMqAdapter, outboxRepository, mongoStream)
	return factory.NewEventFactory(handler)
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
	orderRepository := repository.NewDefaultOrderRepository(connection)
	shoppingCartRepository := repository.NewShoppingCartRedisRepository(redisClient)
	return factory.NewRepositoryFactory(customerRepository, productRepository, shoppingCartRepository, orderRepository)
}
