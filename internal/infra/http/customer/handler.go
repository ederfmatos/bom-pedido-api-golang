package customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/http/customer/get_customer"
	"bom-pedido-api/internal/infra/http/customer/google_auth_customer"
	"bom-pedido-api/internal/infra/http/middlewares"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	server.POST("/v1/auth/google/customer", google_auth_customer.Handle(applicationFactory))
	server.GET("/v1/customers/me", get_customer.Handle(applicationFactory), middlewares.OnlyCustomer)
}
