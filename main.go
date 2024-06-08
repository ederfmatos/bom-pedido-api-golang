package main

import (
	"bom-pedido-api/infra/factory"
	handler2 "bom-pedido-api/presentation/http"
	"bom-pedido-api/presentation/messaging"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	database, err := sql.Open(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer database.Close()

	applicationFactory := factory.NewApplicationFactory(database)

	server := echo.New()
	jaeger := jaegertracing.New(server, nil)
	defer jaeger.Close()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.HTTPErrorHandler = handler2.HandleError

	group := server.Group("/api")
	group.POST("/v1/products", handler2.HandleCreateProduct(applicationFactory))
	group.POST("/v1/auth/google/customer", handler2.HandleGoogleAuthCustomer(applicationFactory))
	group.GET("/v1/customers/me", handler2.HandleGetAuthenticatedCustomer(applicationFactory))
	group.GET("/health", handler2.HandleHealth)

	go applicationFactory.EventHandler.Consume("PRODUCT_CREATED", messaging.HandleCreateProductProjection())

	err = server.Start(":8080")
	server.Logger.Fatal(err)
}
