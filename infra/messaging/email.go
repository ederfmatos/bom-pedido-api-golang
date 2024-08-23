package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"context"
)

func HandleEmailEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopic("SEND_EMAIL", "SEND_EMAIL"), func(ctx context.Context, message *event.MessageEvent) error {
		data := message.Event.Data
		err := factory.EmailGateway.Send(data["to"], data["subject"], data["template"], data)
		return message.AckIfNoError(ctx, err)
	})
}
