package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleProductEvents(factory *factory.ApplicationFactory) {
	go factory.EventDispatcher.Consume("PRODUCTS::CREATE_PRODUCT_PROJECTION", func(event event.Event) error {
		slog.Info("Handling CREATE_PRODUCT_PROJECTION event", "id", event)
		return nil
	})
}
