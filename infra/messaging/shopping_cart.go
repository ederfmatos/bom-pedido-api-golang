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

func handleDeleteShoppingCart(factory *factory.ApplicationFactory) func(message event.MessageEvent) error {
	useCase := delete_shopping_cart.New(factory)
	return func(message event.MessageEvent) error {
		customerId := message.GetEvent().Data["customerId"]
		err := useCase.Execute(context.Background(), delete_shopping_cart.Input{CustomerId: customerId})
		return message.AckIfNoError(err)
	}
}
