package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"github.com/labstack/echo/v4"
)

func HandleGetAuthenticatedCustomer(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := usecase.NewGetCustomerUseCase(factory)
	return func(context echo.Context) error {
		input := usecase.GetCustomerInput{
			Id:      context.Get("customerId").(string),
			Context: context.Request().Context(),
		}
		output, err := useCase.Execute(input)
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
