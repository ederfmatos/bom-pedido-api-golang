package create_product

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/product"
	"bom-pedido-api/internal/infra/http/response"
	"bom-pedido-api/internal/infra/tenant"
	"github.com/labstack/echo/v4"
)

type createProductRequest struct {
	Name        string  `body:"name" json:"name,omitempty"`
	Description string  `body:"description" json:"description,omitempty"`
	CategoryId  string  `body:"categoryId" json:"categoryId,omitempty"`
	Price       float64 `body:"price" json:"price,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	createProductUseCase := product.NewCreateProduct(factory)
	return func(context echo.Context) error {
		var request createProductRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := product.CreateProductInput{
			Name:        request.Name,
			Description: request.Description,
			Price:       request.Price,
			CategoryId:  request.CategoryId,
			TenantId:    context.Get(tenant.Id).(string),
		}
		output, err := createProductUseCase.Execute(context.Request().Context(), input)
		return response.Created(context, output, err)
	}
}
