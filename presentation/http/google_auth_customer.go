package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"github.com/labstack/echo/v4"
)

type GoogleAuthCustomerRequest struct {
	Token string `body:"token"`
}

func HandleGoogleAuthCustomer(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := usecase.NewGoogleAuthenticateCustomerUseCase(factory)
	return func(context echo.Context) error {
		var request GoogleAuthCustomerRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := usecase.GoogleAuthenticateCustomerInput{
			Token:   request.Token,
			Context: context.Request().Context(),
		}
		output, err := useCase.Execute(input)
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
