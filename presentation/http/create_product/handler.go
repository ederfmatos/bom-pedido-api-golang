package create_product

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/product/create_product"
	"github.com/labstack/echo/v4"
)

type createProductRequest struct {
	Name        string  `body:"name" json:"name,omitempty"`
	Description string  `body:"description" json:"description,omitempty"`
	Price       float64 `body:"price" json:"price,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	createProductUseCase := create_product.New(factory)
	return func(context echo.Context) error {
		var request createProductRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := create_product.Input{
			Context:     context.Request().Context(),
			Name:        request.Name,
			Description: request.Description,
			Price:       request.Price,
		}
		output, err := createProductUseCase.Execute(input)
		if err != nil {
			return err
		}
		return context.JSON(201, output)
	}
}
