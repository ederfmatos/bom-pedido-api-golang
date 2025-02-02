package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/usecase/shopping_cart"
	"context"
)

type ShoppingCartHandler struct {
	deleteShoppingCartUseCase *shopping_cart.DeleteShoppingCartUseCase
}

func NewShoppingCartHandlerHandler(deleteShoppingCartUseCase *shopping_cart.DeleteShoppingCartUseCase) *ShoppingCartHandler {
	return &ShoppingCartHandler{deleteShoppingCartUseCase: deleteShoppingCartUseCase}
}

func (h ShoppingCartHandler) DeleteShoppingCart(ctx context.Context, message *event.MessageEvent) error {
	input := shopping_cart.DeleteShoppingCartInput{
		CustomerId: message.Event.Data["customerId"],
	}
	err := h.deleteShoppingCartUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}
