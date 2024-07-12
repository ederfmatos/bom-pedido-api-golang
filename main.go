package main

import (
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/presentation/http"
	"bom-pedido-api/presentation/http/approve_order"
	"bom-pedido-api/presentation/http/create_product"
	"bom-pedido-api/presentation/http/get_customer"
	"bom-pedido-api/presentation/http/health"
	"bom-pedido-api/presentation/http/list_products"
	"bom-pedido-api/presentation/http/middlewares"
	"bom-pedido-api/presentation/http/reject_order"
	"bom-pedido-api/presentation/messaging"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
)

func main() {
	environment := env.LoadEnvironment(".env")
	database, err := sql.Open(environment.DatabaseDriver, environment.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	redisUrl, err := redis.ParseURL(environment.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisUrl)
	defer redisClient.Close()

	applicationFactory := factory.NewApplicationFactory(database, environment, redisClient)
	defer applicationFactory.EventHandler.Close()

	go messaging.HandleEvents(applicationFactory)

	server := echo.New()
	//server.Use(middlewares.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middlewares.AuthenticateMiddleware(applicationFactory))
	server.HTTPErrorHandler = middlewares.HandleError

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	group := server.Group("/api")
	group.POST("/v1/shopping-cart/checkout", http.Handle(applicationFactory), middlewares.LockByCustomerId(applicationFactory))
	group.PATCH("/v1/shopping-cart/items", http.HandleAddItemToShoppingCart(applicationFactory), middlewares.LockByCustomerId(applicationFactory))
	group.POST("/v1/products", create_product.Handle(applicationFactory))
	group.GET("/v1/products", list_products.Handle(applicationFactory))
	group.POST("/v1/auth/google/customer", http.HandleGoogleAuthCustomer(applicationFactory))
	group.GET("/v1/customers/me", get_customer.Handle(applicationFactory))
	group.POST("/v1/orders/:id/approve", approve_order.Handle(applicationFactory), middlewares.LockByParam("id", applicationFactory))
	group.POST("/v1/orders/:id/reject", reject_order.Handle(applicationFactory), middlewares.LockByParam("id", applicationFactory))
	group.GET("/health", health.Handle)

	err = server.Start(fmt.Sprintf(":%s", environment.Port))
	server.Logger.Fatal(err)
}
