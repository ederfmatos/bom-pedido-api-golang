package get_shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/shopping_cart/get_shopping_cart"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := get_shopping_cart.New(factory)
	return func(context echo.Context) error {
		input := get_shopping_cart.Input{
			CustomerId: context.Get(middlewares.CustomerIdParam).(string),
		}
		shoppingCart, err := useCase.Execute(context.Request().Context(), input)
		return response.Ok(context, shoppingCart, err)
	}
}
