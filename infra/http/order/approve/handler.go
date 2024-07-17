package approve

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/approve_order"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := approve_order.New(factory)
	return func(context echo.Context) error {
		input := approve_order.Input{
			OrderId:    context.Param("id"),
			ApprovedBy: context.Get("adminId").(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
