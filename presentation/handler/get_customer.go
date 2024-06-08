package handler

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"github.com/labstack/echo/v4"
)

func HandleGetAuthenticatedCustomer(appFactory *factory.ApplicationFactory) func(context echo.Context) error {
	getCustomerUseCase := usecase.NewGetCustomerUseCase(appFactory)
	return func(context echo.Context) error {
		input := usecase.GetCustomerInput{
			Id:      context.Get("currentUserId").(string),
			Context: context.Request().Context(),
		}
		output, err := getCustomerUseCase.Execute(input)
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
