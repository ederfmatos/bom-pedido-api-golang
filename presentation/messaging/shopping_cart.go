package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/shopping_cart/delete_shopping_cart"
	"context"
)

func HandleShoppingCart(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForQueue("SHOPPING_CART::DELETE_SHOPPING_CART"), handleDeleteShoppingCart(factory))
}

func handleDeleteShoppingCart(factory *factory.ApplicationFactory) func(message event.MessageEvent) error {
	useCase := delete_shopping_cart.New(factory)
	return func(message event.MessageEvent) error {
		var orderCreatedEvent event.OrderEventData
		message.ParseData(&orderCreatedEvent)
		err := useCase.Execute(context.Background(), delete_shopping_cart.Input{CustomerId: orderCreatedEvent.CustomerId})
		return message.AckIfNoError(err)
	}
}
