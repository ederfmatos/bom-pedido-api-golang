package category

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/http/category/create_category"
	"bom-pedido-api/internal/infra/http/middlewares"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	server.POST("/v1/categories", create_category.Handle(applicationFactory), middlewares.OnlyAdmin)
}
