package test

import (
	"bom-pedido-api/internal/infra/config"
	mongo2 "bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

type Container struct {
	MongoClient   *mongo.Client
	mongoDatabase *mongo.Database
	RedisClient   *redis2.Client
	downMongo     func()
	downRedis     func()
	redisURL      string
}

var instance *Container
var ctx = context.Background()
var once sync.Once

func NewContainer() *Container {
	once.Do(func() {
		if instance != nil {
			return
		}
		MongoClient, downMongo := mongoConnection()
		RedisClient, redisURL, downRedis := redisClient()
		instance = &Container{
			MongoClient:   MongoClient,
			mongoDatabase: MongoClient.Database("test"),
			RedisClient:   RedisClient,
			downMongo:     downMongo,
			downRedis:     downRedis,
			redisURL:      redisURL,
		}
	})
	return instance
}

func (c *Container) Down() {
	fmt.Println("Down containers")
	go c.downMongo()
	go c.downRedis()
}

func mongoConnection() (*mongo.Client, func()) {
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	failOnError(err)
	endpoint, err := mongodbContainer.Endpoint(context.Background(), "")
	failOnError(err)
	clientOptions := options.Client().ApplyURI("mongodb://" + endpoint)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	failOnError(err)
	return mongoClient, func() {
		_ = mongodbContainer.Terminate(ctx)
		_ = mongoClient.Disconnect(ctx)
	}
}

func redisClient() (*redis2.Client, string, func()) {
	redisContainer, err := redis.Run(ctx, "docker.io/redis:7")
	failOnError(err)
	endpoint, err := redisContainer.Endpoint(ctx, "")
	failOnError(err)
	redisUrl, err := redis2.ParseURL("redis://" + endpoint)
	failOnError(err)
	redisClient := redis2.NewClient(redisUrl)
	return redisClient, "redis://" + endpoint, func() {
		_ = redisClient.Close()
		_ = redisContainer.Terminate(ctx)
	}
}

func (c *Container) MongoDatabase() *mongo2.Database {
	return mongo2.NewDatabase(c.mongoDatabase)
}

func (c *Container) GetEnvironment() *config.Environment {
	return &config.Environment{
		RedisUrl:                      c.redisURL,
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
