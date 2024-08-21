package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Container struct {
	Database      *sql.DB
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
	RedisClient   *redis2.Client
	downDatabase  func()
	downMongo     func()
	downRedis     func()
}

var instance *Container
var ctx = context.TODO()

func init() {
	if instance != nil {
		return
	}
	fmt.Println("Criando nova instancia do instance")
	MongoClient, downMongo := mongoConnection()
	RedisClient, downRedis := redisClient()
	Database, downDatabase := databaseConnection()
	instance = &Container{
		Database:      Database,
		MongoClient:   MongoClient,
		MongoDatabase: MongoClient.Database("test"),
		RedisClient:   RedisClient,
		downDatabase:  downDatabase,
		downMongo:     downMongo,
		downRedis:     downRedis,
	}
}

func NewContainer() *Container {
	return instance
}

func (c *Container) Down() {
	fmt.Println("Down containers")
	go c.downMongo()
	go c.downRedis()
	go c.downDatabase()
}

func mongoConnection() (*mongo.Client, func()) {
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	endpoint, err := mongodbContainer.Endpoint(context.Background(), "")
	failOnError(err)
	clientOptions := options.Client().ApplyURI("mongodb://" + endpoint)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	failOnError(err)
	return mongoClient, func() {
		mongodbContainer.Terminate(ctx)
		mongoClient.Disconnect(ctx)
	}
}

func redisClient() (*redis2.Client, func()) {
	redisContainer, err := redis.Run(ctx, "docker.io/redis:7")
	failOnError(err)
	endpoint, err := redisContainer.Endpoint(ctx, "")
	failOnError(err)
	redisUrl, err := redis2.ParseURL("redis://" + endpoint)
	failOnError(err)
	redisClient := redis2.NewClient(redisUrl)
	return redisClient, func() {
		redisClient.Close()
		redisContainer.Terminate(ctx)
	}
}

func databaseConnection() (*sql.DB, func()) {
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	failOnError(err)
	connectionString := postgresContainer.MustConnectionString(ctx, "sslmode=disable")
	database, err := sql.Open("postgres", connectionString)
	failOnError(err)
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS customers
		(
			id           VARCHAR(36)                NOT NULL PRIMARY KEY,
			name         VARCHAR(255)               NOT NULL,
			email        VARCHAR(255)               NOT NULL UNIQUE,
			phone_number VARCHAR(11)                UNIQUE,
			status       VARCHAR(20)				NOT NULL,
			created_at   TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
    `)
	failOnError(err)
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS products
		(
			id          VARCHAR(36)                            NOT NULL PRIMARY KEY,
			name        VARCHAR(255)                           NOT NULL UNIQUE,
			description TEXT,
			price       DECIMAL(6, 2)                          NOT NULL,
			status       VARCHAR(20)						   NOT NULL,
			created_at  TIMESTAMP                              NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
    `)
	failOnError(err)
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS orders
		(
			id                VARCHAR(36) NOT NULL PRIMARY KEY,
			code              SERIAL      NOT NULL UNIQUE,
			customer_id       VARCHAR(36) NOT NULL,
			payment_method    VARCHAR(30) NOT NULL,
			payment_mode      VARCHAR(30) NOT NULL,
			delivery_mode     VARCHAR(30) NOT NULL,
			status            VARCHAR(30) NOT NULL,
			credit_card_token VARCHAR(255),
			payback           DECIMAL(6, 2),
			delivery_time     TIMESTAMP   NOT NULL,
			created_at        TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
    `)
	failOnError(err)
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS order_items
		(
			id          VARCHAR(36) NOT NULL PRIMARY KEY,
			order_id    VARCHAR(36) NOT NULL,
			product_id  VARCHAR(36) NOT NULL,
			status      VARCHAR(30) NOT NULL,
			quantity    NUMERIC     NOT NULL,
			observation TEXT,
			price       DECIMAL(6, 2),
			created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id)
		);
    `)
	failOnError(err)
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS order_history
		(
			id         SERIAL      NOT NULL PRIMARY KEY,
			order_id   VARCHAR(36) NOT NULL,
			changed_by VARCHAR(36) NOT NULL,
			changed_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status     VARCHAR(30) NOT NULL,
			data       TEXT,
			CONSTRAINT fk_order_history_order FOREIGN KEY (order_id) REFERENCES orders (id)
		);
    `)
	failOnError(err)
	return database, func() {
		database.Close()
		postgresContainer.Terminate(ctx)
	}
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
