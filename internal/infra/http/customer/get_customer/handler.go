package get_customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/customer/get_customer"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := get_customer.New(factory)
	return func(context echo.Context) error {
		input := get_customer.Input{
			Id: context.Get(middlewares.CustomerIdParam).(string),
		}
		output, err := useCase.Execute(context.Request().Context(), input)
		return response.Ok(context, output, err)
	}
}
