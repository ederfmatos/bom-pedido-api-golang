package main

import (
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/presentation/http"
	middleware2 "bom-pedido-api/presentation/http/middleware"
	"bom-pedido-api/presentation/messaging"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"
)

func main() {
	environment := env.LoadEnvironment(".environment")
	database, err := sql.Open(environment.DatabaseDriver, environment.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	applicationFactory := factory.NewApplicationFactory(database, environment)

	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middleware2.AuthenticateMiddleware(applicationFactory))
	server.HTTPErrorHandler = middleware2.HandleError

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	group := server.Group("/api")
	group.PATCH("/v1/shopping-cart/items", http.HandleAddItemToShoppingCart(applicationFactory))
	group.POST("/v1/products", http.HandleCreateProduct(applicationFactory))
	group.GET("/v1/products", http.HandleListProducts(applicationFactory))
	group.POST("/v1/auth/google/customer", http.HandleGoogleAuthCustomer(applicationFactory))
	group.GET("/v1/customers/me", http.HandleGetAuthenticatedCustomer(applicationFactory))
	group.GET("/health", http.HandleHealth)

	eventHandler := applicationFactory.EventHandler
	go eventHandler.Consume("CREATE_PRODUCT_PROJECTION", messaging.HandleCreateProductProjection())

	err = server.Start(":8080")
	server.Logger.Fatal(err)
}
