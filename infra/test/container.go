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
	"sync"
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
var once sync.Once

func NewContainer() *Container {
	once.Do(func() {
		if instance != nil {
			return
		}
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
	})
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
		CREATE TABLE merchants
		(
			id           VARCHAR(36) PRIMARY KEY,
			name         VARCHAR(255) NOT NULL,
			email        VARCHAR(255) NOT NULL,
			phone_number VARCHAR(11)  NOT NULL,
			tenant_id    VARCHAR(36)  NOT NULL UNIQUE,
			domain       VARCHAR(20)  NOT NULL,
			status       VARCHAR(30)  NOT NULL DEFAULT 'ACTIVE',
			created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uk_merchant_email UNIQUE (email, tenant_id),
			CONSTRAINT uk_merchant_phone_number UNIQUE (phone_number, tenant_id)
		);
		CREATE INDEX ids_merchants_domain ON merchants (domain);
		CREATE INDEX ids_merchants_tenant_id ON merchants (tenant_id);
		
		CREATE TABLE merchant_address
		(
			id           SERIAL PRIMARY KEY,
			merchant_id  VARCHAR(36)  NOT NULL,
			street       VARCHAR(255) NOT NULL,
			number       VARCHAR(20),
			neighborhood VARCHAR(255),
			postal_code  VARCHAR(8),
			city         VARCHAR(100) NOT NULL,
			state        VARCHAR(2)   NOT NULL,
			created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_merchant_address_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
		);
		
		CREATE INDEX ids_merchant_address_merchant_id ON merchant_address (merchant_id);
		
		CREATE TABLE merchant_opening_hour
		(
			id           SERIAL PRIMARY KEY,
			merchant_id  VARCHAR(36) NOT NULL,
			day_of_week  NUMERIC(2)  NOT NULL,
			initial_time TIME        NOT NULL,
			final_time   TIME        NOT NULL,
			created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_merchant_opening_hour_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
		);
		CREATE INDEX ids_merchant_opening_hour_merchant_id ON merchant_opening_hour (merchant_id);
		
		CREATE TABLE merchant_payment_gateway_configs
		(
			id          SERIAL      NOT NULL PRIMARY KEY,
			merchant_id VARCHAR(36) NOT NULL,
			gateway     VARCHAR(50) NOT NULL,
			credentials TEXT        NOT NULL,
			CONSTRAINT fk_merchant_payment_gateway_configs_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
		);
		CREATE INDEX ids_merchant_payment_gateway_configs_merchant_id ON merchant_payment_gateway_configs (merchant_id);
		
		CREATE TABLE admins
		(
			id          VARCHAR(36) PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			email       VARCHAR(255) NOT NULL,
			merchant_id VARCHAR(36)  NOT NULL,
			created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_admins_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id),
			CONSTRAINT uk_admin_email UNIQUE (email, merchant_id)
		);
		CREATE INDEX ids_admin_email ON admins (email);
		
		CREATE TABLE customers
		(
			id           VARCHAR(36)  NOT NULL PRIMARY KEY,
			name         VARCHAR(255) NOT NULL,
			email        VARCHAR(255) NOT NULL,
			phone_number VARCHAR(11),
			status       VARCHAR(20)  NOT NULL,
			tenant_id    VARCHAR(30)  NOT NULL,
			created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uk_customer_email UNIQUE (email, tenant_id),
			CONSTRAINT uk_customer_phone_number UNIQUE (phone_number, tenant_id)
		);
		CREATE INDEX ids_customer_email_tenant ON customers (email, tenant_id);
		
		CREATE TABLE products
		(
			id          VARCHAR(36)   NOT NULL PRIMARY KEY,
			name        VARCHAR(255)  NOT NULL,
			description TEXT,
			price       DECIMAL(6, 2) NOT NULL,
			status      VARCHAR(20)   NOT NULL,
			tenant_id   VARCHAR(30)   NOT NULL,
			created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		CONSTRAINT uk_product_name UNIQUE (name, tenant_id)
		);
		CREATE INDEX ids_product_tenant ON products (tenant_id);
		
		CREATE TABLE orders
		(
			id                VARCHAR(36)   NOT NULL PRIMARY KEY,
			code              SERIAL        NOT NULL,
			customer_id       VARCHAR(36)   NOT NULL,
			payment_method    VARCHAR(30)   NOT NULL,
			payment_mode      VARCHAR(30)   NOT NULL,
			delivery_mode     VARCHAR(30)   NOT NULL,
			status            VARCHAR(30)   NOT NULL,
			credit_card_token VARCHAR(255),
			payback           DECIMAL(6, 2),
			amount            DECIMAL(6, 2) NOT NULL,
			delivery_time     TIMESTAMP     NOT NULL,
			merchant_id       VARCHAR(36)   NOT NULL,
			created_at        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_orders_customer FOREIGN KEY (customer_id) REFERENCES customers (id)
		);
		CREATE INDEX ids_order_customer ON orders (customer_id, merchant_id);
		
		CREATE TABLE order_items
		(
			id          VARCHAR(36) NOT NULL PRIMARY KEY,
			order_id    VARCHAR(36) NOT NULL,
			product_id  VARCHAR(36) NOT NULL,
			status      VARCHAR(30) NOT NULL,
			quantity    NUMERIC     NOT NULL,
			observation TEXT,
			price       DECIMAL(6, 2),
			created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id),
			CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products (id)
		);
		
		CREATE TABLE order_history
		(
			id         SERIAL      NOT NULL PRIMARY KEY,
			order_id   VARCHAR(36) NOT NULL,
			changed_by VARCHAR(36) NOT NULL,
			changed_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status     VARCHAR(30) NOT NULL,
			data       TEXT,
			CONSTRAINT fk_order_history_order FOREIGN KEY (order_id) REFERENCES orders (id)
		);
		
		CREATE TYPE transaction_type AS ENUM ('PIX', 'CREDIT_CARD');
		
		CREATE TABLE transactions
		(
			id         VARCHAR(36) PRIMARY KEY NOT NULL,
			order_id   VARCHAR(36)             NOT NULL,
    		payment_id VARCHAR(100)            NOT NULL,
			amount     NUMERIC(6, 2)           NOT NULL,
			status     VARCHAR(20)             NOT NULL,
			created_at TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
			type       transaction_type        NOT NULL,
			CONSTRAINT fk_transactions_order_id FOREIGN KEY (order_id) REFERENCES orders (id)
		);
		CREATE INDEX idx_transactions_order ON transactions (order_id);
		
		CREATE TABLE pix_transactions
		(
			qr_code         TEXT        NOT NULL,
			qr_code_link    TEXT        NOT NULL,
			payment_gateway VARCHAR(50) NOT NULL
		) inherits (transactions);
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
