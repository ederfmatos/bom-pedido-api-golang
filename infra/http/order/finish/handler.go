package finish

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/finish_order"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := finish_order.New(factory)
	return func(context echo.Context) error {
		input := finish_order.Input{
			OrderId:    context.Param("id"),
			FinishedBy: context.Get("adminId").(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		if err != nil {
			return err
		}
		return context.NoContent(204)
	}
}
