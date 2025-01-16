package google_auth_customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/customer"
	"bom-pedido-api/internal/infra/http/response"
	"bom-pedido-api/internal/infra/tenant"
	"github.com/labstack/echo/v4"
)

type GoogleAuthCustomerRequest struct {
	Token string `body:"token"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := customer.NewGoogleAuthenticateCustomer(factory)
	return func(context echo.Context) error {
		var request GoogleAuthCustomerRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := customer.GoogleAuthenticateCustomerInput{
			Token:    request.Token,
			TenantId: context.Get(tenant.Id).(string),
		}
		output, err := useCase.Execute(context.Request().Context(), input)
		return response.Ok(context, output, err)
	}
}
