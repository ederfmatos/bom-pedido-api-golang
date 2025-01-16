package cancel

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/order"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

type cancelOrderRequest struct {
	Reason string `body:"reason" json:"reason,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := order.NewCancelOrder(factory)
	return func(context echo.Context) error {
		var request cancelOrderRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := order.CancelOrderInput{
			OrderId:     context.Param("id"),
			CancelledBy: context.Get(middlewares.AdminIdParam).(string),
			Reason:      request.Reason,
		}
		err = useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
