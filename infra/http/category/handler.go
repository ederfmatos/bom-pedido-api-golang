package category

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/category/create_category"
	"bom-pedido-api/infra/http/middlewares"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	server.POST("/v1/categories", create_category.Handle(applicationFactory), middlewares.OnlyAdmin)
}
