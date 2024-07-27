package health

import (
	"bom-pedido-api/infra/config"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_Health(t *testing.T) {
	instance := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	database, closeDatabase := DatabaseConnection()
	mongoClient, closeMongoConnection := MongoConnection()
	redisClient, closeRedisConnection := RedisClient()
	defer func() {
		go closeDatabase()
		go closeMongoConnection()
		go closeRedisConnection()
	}()

	err := Handle(database, redisClient, mongoClient)(instance.NewContext(request, response))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)

	var output Output
	_ = json.NewDecoder(response.Body).Decode(&output)
	assert.Equal(t, true, output.Ok)
}

func MongoConnection() (*mongo.Client, func()) {
	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	connectionString, err := mongodbContainer.ConnectionString(context.Background())
	mongoClient := config.Mongo(connectionString)
	failOnError(err)
	return mongoClient, func() {
		mongodbContainer.Terminate(ctx)
		mongoClient.Disconnect(ctx)
	}
}

func RedisClient() (*redis2.Client, func()) {
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx, "docker.io/redis:7")
	failOnError(err)
	connectionString, err := redisContainer.ConnectionString(ctx)
	failOnError(err)
	redisClient := config.Redis(connectionString)
	return redisClient, func() {
		redisClient.Close()
		redisContainer.Terminate(ctx)
	}
}

func DatabaseConnection() (*sql.DB, func()) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	failOnError(err)
	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	failOnError(err)
	database := config.Database("postgres", connectionString)
	return database, func() {
		database.Close()
		postgresContainer.Terminate(ctx)
	}
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}
