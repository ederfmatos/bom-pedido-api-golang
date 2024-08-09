package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/shopping_cart/delete_shopping_cart"
	"context"
)

func HandleShoppingCart(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopic("ORDER_CREATED", "DELETE_SHOPPING_CART"), handleDeleteShoppingCart(factory))
}

func handleDeleteShoppingCart(factory *factory.ApplicationFactory) func(context.Context, *event.MessageEvent) error {
	useCase := delete_shopping_cart.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		customerId := message.Event.Data["customerId"]
		err := useCase.Execute(ctx, delete_shopping_cart.Input{CustomerId: customerId})
		return message.AckIfNoError(ctx, err)
	}
}
