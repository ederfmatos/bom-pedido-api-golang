package finish

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/order/finish_order"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := finish_order.New(factory)
	return func(context echo.Context) error {
		input := finish_order.Input{
			OrderId:    context.Param("id"),
			FinishedBy: context.Get(middlewares.AdminIdParam).(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
