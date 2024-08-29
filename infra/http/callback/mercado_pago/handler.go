package mercado_pago

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"github.com/labstack/echo/v4"
)

const gateway = "MERCADO_PAGO"

type callbackRequest struct {
	Action string `body:"action" json:"action,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request callbackRequest
		err := c.Bind(&request)
		if err != nil || request.Action != "payment.updated" {
			return err
		}
		orderId := c.Param("orderId")
		callbackEvent := event.NewPaymentCallbackReceived(gateway, orderId, request.Action)
		return factory.EventEmitter.Emit(c.Request().Context(), callbackEvent)
	}
}
