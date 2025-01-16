package create_category

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/category"
	"bom-pedido-api/internal/infra/http/response"
	"bom-pedido-api/internal/infra/tenant"
	"github.com/labstack/echo/v4"
)

type requestBody struct {
	Name        string `body:"name" json:"name,omitempty"`
	Description string `body:"description" json:"description,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	createProductUseCase := category.NewCreateCategory(factory)
	return func(context echo.Context) error {
		var request requestBody
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := category.CreateCategoryInput{
			Name:        request.Name,
			Description: request.Description,
			TenantId:    context.Get(tenant.Id).(string),
		}
		err = createProductUseCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
