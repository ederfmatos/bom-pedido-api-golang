package woovi

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"github.com/labstack/echo/v4"
)

const gateway = "WOOVI"

type callbackRequest struct {
	Event  string `body:"event"`
	Charge struct {
		CorrelationID string `body:"correlationID"`
	} `body:"charge"`
}

func Handle(factory *factory.ApplicationFactory) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request callbackRequest
		err := c.Bind(&request)
		if err != nil || request.Event != "OPENPIX:CHARGE_COMPLETED" {
			return err
		}
		orderId := request.Charge.CorrelationID
		if orderId == "" {
			return nil
		}
		callbackEvent := event.NewPaymentCallbackReceived(gateway, orderId, "PAYMENT_COMPLETED")
		return factory.EventEmitter.Emit(c.Request().Context(), callbackEvent)
	}
}
