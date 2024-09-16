package create_category

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/product/category/create_category"
	"bom-pedido-api/infra/http/response"
	"bom-pedido-api/infra/tenant"
	"github.com/labstack/echo/v4"
)

type requestBody struct {
	Name        string `body:"name" json:"name,omitempty"`
	Description string `body:"description" json:"description,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	createProductUseCase := create_category.New(factory)
	return func(context echo.Context) error {
		var request requestBody
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := create_category.Input{
			Name:        request.Name,
			Description: request.Description,
			TenantId:    context.Get(tenant.Id).(string),
		}
		err = createProductUseCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
