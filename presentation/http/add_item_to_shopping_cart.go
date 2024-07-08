package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/shopping_cart/add_item_to_shopping_cart"
	"github.com/labstack/echo/v4"
)

type addItemToShoppingCartRequest struct {
	ProductId   string `body:"productId" json:"productId,omitempty"`
	Quantity    int    `body:"quantity" json:"quantity,omitempty"`
	Observation string `body:"observation" json:"observation,omitempty"`
}

func HandleAddItemToShoppingCart(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := add_item_to_shopping_cart.New(factory)
	return func(context echo.Context) error {
		var request addItemToShoppingCartRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := add_item_to_shopping_cart.Input{
			Context:     context.Request().Context(),
			CustomerId:  context.Get("customerId").(string),
			ProductId:   request.ProductId,
			Quantity:    request.Quantity,
			Observation: request.Observation,
		}
		err = useCase.Execute(input)
		if err != nil {
			return err
		}
		return context.NoContent(204)
	}
}
