package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	echoredoc "github.com/mvrilo/go-redoc/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

func main() {
	server := echo.New()

	doc := redoc.Redoc{
		Title:       "Bom Pedido API",
		Description: "API para gerenciamento de pedidos, delivery e fisicos.",
		SpecFile:    ".docs/openapi.json",
		SpecPath:    "openapi.json",
		DocsPath:    "/docs",
	}
	server.Use(echoredoc.New(doc))
	server.GET("/swagger.json", func(c echo.Context) error {
		return c.File(".docs/openapi.json")
	})
	server.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{".docs/openapi.json"}
	}))

	err := server.Start(":8080")
	log.Fatal(err)
}
