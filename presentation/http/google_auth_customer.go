package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/customer/google_authenticate_customer"
	"github.com/labstack/echo/v4"
)

type GoogleAuthCustomerRequest struct {
	Token string `body:"token"`
}

func HandleGoogleAuthCustomer(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := google_authenticate_customer.New(factory)
	return func(context echo.Context) error {
		var request GoogleAuthCustomerRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := google_authenticate_customer.Input{
			Token: request.Token,
		}
		output, err := useCase.Execute(context.Request().Context(), input)
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
