package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleProductEvents(factory *factory.ApplicationFactory) {
	factory.EventDispatcher.Consume(event.OptionsForQueue("PRODUCTS::CREATE_PRODUCT_PROJECTION"), func(message event.MessageEvent) error {
		slog.Info("Received product message", "message", message.AsEvent())
		return message.Ack()
	})
}
