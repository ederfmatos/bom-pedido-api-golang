package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
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

func RedisClient() (*redis2.Client, func()) {
	ctx := context.Background()
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

func DatabaseConnection() (*sql.DB, func()) {
	ctx := context.Background()
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	failOnError(err)
	endpoint, err := mysqlContainer.Endpoint(ctx, "")
	failOnError(err)
	database, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s)/test", endpoint))
	failOnError(err)
	return database, func() {
		database.Close()
		mysqlContainer.Terminate(ctx)
	}
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}
