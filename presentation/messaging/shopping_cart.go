package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleShoppingCart(factory *factory.ApplicationFactory) {
	go factory.EventDispatcher.Consume(event.OptionsForQueue("SHOPPING_CART::DELETE_SHOPPING_CART"), func(message event.MessageEvent) error {
		slog.Info("Handling DELETE_SHOPPING_CART event", "id", message.AsEvent())
		return message.Ack()
	})
}
