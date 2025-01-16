package mark_delivering

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/order"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := order.NewMarkOrderDeliveringUseCase(factory)
	return func(context echo.Context) error {
		input := order.MarkOrderDeliveringInput{
			OrderId: context.Param("id"),
			By:      context.Get(middlewares.AdminIdParam).(string),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
