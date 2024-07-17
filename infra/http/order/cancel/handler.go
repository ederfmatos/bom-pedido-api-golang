package cancel

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/cancel_order"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

type cancelOrderRequest struct {
	Reason string `body:"reason" json:"reason,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := cancel_order.New(factory)
	return func(context echo.Context) error {
		var request cancelOrderRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := cancel_order.Input{
			OrderId:     context.Param("id"),
			CancelledBy: context.Get("adminId").(string),
			Reason:      request.Reason,
		}
		err = useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
