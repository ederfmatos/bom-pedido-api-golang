package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	echoredoc "github.com/mvrilo/go-redoc/echo"
)

func RedocDocumentation() echo.MiddlewareFunc {
	doc := redoc.Redoc{
		Title:       "Bom Pedido API",
		Description: "API para gerenciamento de pedidos, delivery e fisicos.",
		SpecFile:    ".docs/openapi.json",
		SpecPath:    "openapi.json",
		DocsPath:    "/docs",
	}
	return echoredoc.New(doc)
}
