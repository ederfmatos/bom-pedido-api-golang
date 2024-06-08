package handler

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"github.com/labstack/echo/v4"
)

type createProductRequest struct {
	Name        string  `body:"name" json:"name,omitempty"`
	Description string  `body:"description" json:"description,omitempty"`
	Price       float64 `body:"price" json:"price,omitempty"`
}

func HandleCreateProduct(appFactory *factory.ApplicationFactory) func(context echo.Context) error {
	createProductUseCase := usecase.NewCreateProductUseCase(appFactory)
	return func(context echo.Context) error {
		var request createProductRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := usecase.CreateProductInput{
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
