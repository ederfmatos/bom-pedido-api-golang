package send_authentication_magic_link

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/admin/send_authentication_magic_link"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

type requestBody struct {
	Email string `body:"email"`
}

func Handle(baseURL string, factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := send_authentication_magic_link.New(baseURL, factory)
	return func(context echo.Context) error {
		var request requestBody
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := send_authentication_magic_link.Input{Email: request.Email}
		err = useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
