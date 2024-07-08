package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleShoppingCart(factory *factory.ApplicationFactory) {
	go factory.EventDispatcher.Consume("SHOPPING_CART::DELETE_SHOPPING_CART", func(event event.Event) error {
		slog.Info("Handling DELETE_SHOPPING_CART event", "id", event.Data)
		return nil
	})
}
