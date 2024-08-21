package clone

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/clone_order"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := clone_order.New(factory)
	return func(context echo.Context) error {
		input := clone_order.Input{
			OrderId: context.Param("id"),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
