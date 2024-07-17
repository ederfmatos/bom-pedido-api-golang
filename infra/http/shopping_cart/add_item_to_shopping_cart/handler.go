package add_item_to_shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/shopping_cart/add_item_to_shopping_cart"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

type addItemToShoppingCartRequest struct {
	ProductId   string `body:"productId" json:"productId,omitempty"`
	Quantity    int    `body:"quantity" json:"quantity,omitempty"`
	Observation string `body:"observation" json:"observation,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := add_item_to_shopping_cart.New(factory)
	return func(context echo.Context) error {
		var request addItemToShoppingCartRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := add_item_to_shopping_cart.Input{
			CustomerId:  context.Get("customerId").(string),
			ProductId:   request.ProductId,
			Quantity:    request.Quantity,
			Observation: request.Observation,
		}
		err = useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
