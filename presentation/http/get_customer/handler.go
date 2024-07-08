package get_customer

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/customer/get_customer"
	"github.com/labstack/echo/v4"
)

func HandleGetAuthenticatedCustomer(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := get_customer.New(factory)
	return func(context echo.Context) error {
		input := get_customer.Input{
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
