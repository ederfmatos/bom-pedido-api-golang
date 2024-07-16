package mark_in_progress

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/mark_order_in_progress"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := mark_order_in_progress.New(factory)
	return func(context echo.Context) error {
		input := mark_order_in_progress.Input{
			OrderId: context.Param("id"),
			By:      context.Get("adminId").(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		if err != nil {
			return err
		}
		return context.NoContent(204)
	}
}
