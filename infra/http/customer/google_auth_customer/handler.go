package google_auth_customer

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/customer/google_authenticate_customer"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

type GoogleAuthCustomerRequest struct {
	Token string `body:"token"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
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
		return response.Ok(context, output, err)
	}
}
