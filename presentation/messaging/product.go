package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleProductEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForQueue("PRODUCTS::CREATE_PRODUCT_PROJECTION"), func(message event.MessageEvent) error {
		slog.Info("Received product message")
		return message.Ack()
	})
}
