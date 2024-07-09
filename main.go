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
	environment := env.LoadEnvironment(".env")
	database, err := sql.Open(environment.DatabaseDriver, environment.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	applicationFactory := factory.NewApplicationFactory(database, environment)
	defer applicationFactory.EventHandler.Close()

	go messaging.HandleEvents(applicationFactory)

	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middleware2.AuthenticateMiddleware(applicationFactory))
	server.HTTPErrorHandler = middleware2.HandleError

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	group := server.Group("/api")
	group.POST("/v1/shopping-cart/checkout", http.Handle(applicationFactory))
	group.PATCH("/v1/shopping-cart/items", http.HandleAddItemToShoppingCart(applicationFactory))
	group.POST("/v1/products", create_product.Handle(applicationFactory))
	group.GET("/v1/products", list_products.Handle(applicationFactory))
	group.POST("/v1/auth/google/customer", http.HandleGoogleAuthCustomer(applicationFactory))
	group.GET("/v1/customers/me", get_customer.Handle(applicationFactory))
	group.POST("/v1/orders/:id/approve", approve_order.Handle(applicationFactory))
	group.GET("/health", health.Handle)

	err = server.Start(":8080")
	server.Logger.Fatal(err)
}
