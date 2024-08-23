package admin

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/http/admin/send_authentication_magic_link"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory, environment *config.Environment) {
	server.POST("/v1/auth/admin", send_authentication_magic_link.Handle(environment.AdminMagicLinkBaseUrl, applicationFactory))
}
