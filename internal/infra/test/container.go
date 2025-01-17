package test

import (
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/pkg/mongo"
	"bom-pedido-api/pkg/testcontainer"
	"context"
	"log"
	"sync"
)

type Container struct {
	*testcontainer.RedisContainer
	*testcontainer.MongoContainer
}

var instance *Container
var once sync.Once

func NewContainer() *Container {
	once.Do(func() {
		if instance != nil {
			return
		}

		mongoContainer, err := testcontainer.NewMongoContainer(context.Background())
		failOnError(err)

		redisContainer, err := testcontainer.NewRedisContainer(context.Background())
		failOnError(err)

		instance = &Container{
			RedisContainer: redisContainer,
			MongoContainer: mongoContainer,
		}
	})
	return instance
}

func (c *Container) Down() {
	go c.RedisContainer.Shutdown(context.Background())
	go c.MongoContainer.Shutdown(context.Background())
}

func (c *Container) MongoDatabase() *mongo.Database {
	return mongo.NewDatabase(c.MongoContainer.MongoClient.Database("test"))
}

func (c *Container) GetEnvironment() *config.Environment {
	return &config.Environment{
		RedisUrl:                      c.RedisContainer.Address,
		JwePrivateKeyPath:             "",
		RabbitMqServer:                "",
		GoogleAuthUrl:                 "",
		MongoUrl:                      "",
		MongoDatabaseName:             "",
		MongoOutboxCollectionName:     "",
		Port:                          "",
		OpenTelemetryEndpointExporter: "",
		MessagingStrategy:             "",
		AdminMagicLinkBaseUrl:         "",
		EmailFrom:                     "",
		ResendMailKey:                 "",
		PixPaymentGateway:             config.PixPaymentGatewayEnv{},
	}
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
