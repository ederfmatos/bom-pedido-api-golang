package callback

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/callback/mercado_pago"
	"bom-pedido-api/infra/http/callback/woovi"
	"github.com/labstack/echo/v4"
)

func ConfigureCallbackRoutes(server *echo.Group, factory *factory.ApplicationFactory) {
	callbackRoutes := server.Group("/v1/payments/callback")
	callbackRoutes.POST("/MERCADO_PAGO/:orderId", mercado_pago.Handle(factory))
	callbackRoutes.POST("/WOOVI", woovi.Handle(factory))
}
