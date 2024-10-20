package products

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/products/create_product"
	"bom-pedido-api/internal/infra/http/products/list_products"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	server.POST("/v1/products", create_product.Handle(applicationFactory), middlewares.OnlyAdmin)
	server.GET("/v1/products", list_products.Handle(applicationFactory))
}
