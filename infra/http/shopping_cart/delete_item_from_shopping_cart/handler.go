package delete_item_from_shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/shopping_cart/delete_shopping_cart_item"
	"bom-pedido-api/infra/http/middlewares"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := delete_shopping_cart_item.New(factory)
	return func(context echo.Context) error {
		input := delete_shopping_cart_item.Input{
			CustomerId: context.Get(middlewares.CustomerIdParam).(string),
			ItemId:     context.Param("id"),
		}
		err := useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
