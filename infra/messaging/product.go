package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"context"
	"log/slog"
)

func HandleProductEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopic("PRODUCT_CREATED", "CREATE_PRODUCT_PROJECTION"), func(ctx context.Context, message *event.MessageEvent) error {
		slog.Info("Received product message", "productId", message.Event.Data["productId"])
		return message.Ack(ctx)
	})
}
