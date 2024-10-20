package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"context"
)

func HandleEmailEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("SEND_EMAIL", event.SendEmail), func(ctx context.Context, message *event.MessageEvent) error {
		data := message.Event.Data
		err := factory.EmailGateway.Send(data["to"], data["subject"], data["template"], data)
		return message.AckIfNoError(ctx, err)
	})
}
