package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/transaction"
	"context"
)

func HandleTransactionCallback(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("PAY_PIX_TRANSACTION", event.PaymentCallbackReceived), handlePayPixTransaction(factory))
}

func handlePayPixTransaction(factory *factory.ApplicationFactory) func(context.Context, *event.MessageEvent) error {
	useCase := transaction.NewPayPixTransaction(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := transaction.PayPixTransactionInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
