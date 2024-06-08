package main

import (
	"bom-pedido-api/infra/factory"
	handler2 "bom-pedido-api/presentation/handler"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	_ = os.Setenv("GOOGLE_AUTH_URL", "https://www.googleapis.com/oauth2/v2/userinfo")

	database, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer database.Close()

	applicationFactory := factory.NewApplicationFactory(database)

	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())

	server.POST("/v1/products", handler2.HandleCreateProduct(applicationFactory))
	server.POST("/v1/auth/google/customer", handler2.HandleGoogleAuthCustomer(applicationFactory))
	server.GET("/v1/customers/me", handler2.HandleGetAuthenticatedCustomer(applicationFactory))
	server.GET("/health", handler2.HandleHealth)

	err = server.Start(":8080")
	server.Logger.Fatal(err)
}
