package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"log/slog"
)

func HandleProductEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopic("PRODUCT_CREATED", "CREATE_PRODUCT_PROJECTION"), func(message *event.MessageEvent) error {
		slog.Info("Received product message", "productId", message.GetEvent().Data["productId"])
		return message.Ack()
	})
}
