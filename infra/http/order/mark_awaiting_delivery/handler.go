package mark_awaiting_delivery

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/mark_order_awaiting_delivery"
	"bom-pedido-api/infra/http/middlewares"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := mark_order_awaiting_delivery.New(factory)
	return func(context echo.Context) error {
		input := mark_order_awaiting_delivery.Input{
			OrderId: context.Param("id"),
			By:      context.Get(middlewares.AdminIdParam).(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
