package approve_order

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/approve_order"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := approve_order.New(factory)
	return func(context echo.Context) error {
		input := approve_order.Input{
			Context:    context.Request().Context(),
			OrderId:    context.Param("id"),
			ApprovedBy: context.Get("adminId").(string),
		}
		err := useCase.Execute(input)
		if err != nil {
			return err
		}
		return context.NoContent(204)
	}
}
