package get_customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/customer"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := customer.NewGetCustomer(factory)
	return func(context echo.Context) error {
		input := customer.GetCustomerInput{
			Id: context.Get(middlewares.CustomerIdParam).(string),
		}
		output, err := useCase.Execute(context.Request().Context(), input)
		return response.Ok(context, output, err)
	}
}
