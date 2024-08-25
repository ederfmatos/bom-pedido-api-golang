package callback

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/callback/mercado_pago"
	"github.com/labstack/echo/v4"
)

func ConfigureCallbackRoutes(server *echo.Group, factory *factory.ApplicationFactory) {
	callbackRoutes := server.Group("/v1/payments/callback")
	callbackRoutes.POST("/mercado_pago/:orderId", mercado_pago.Handle(factory))
}
