package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/shopping_cart"
	"context"
)

func HandleShoppingCart(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("DELETE_SHOPPING_CART", "ORDER_CREATED"), handleDeleteShoppingCart(factory))
}

func handleDeleteShoppingCart(factory *factory.ApplicationFactory) func(context.Context, *event.MessageEvent) error {
	useCase := shopping_cart.NewDeleteShoppingCart(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		customerId := message.Event.Data["customerId"]
		err := useCase.Execute(ctx, shopping_cart.DeleteShoppingCartInput{CustomerId: customerId})
		return message.AckIfNoError(ctx, err)
	}
}
