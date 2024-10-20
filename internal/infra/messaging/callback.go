package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/transaction/pay_pix_transaction"
	"context"
)

func HandleTransactionCallback(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("PAY_PIX_TRANSACTION", event.PaymentCallbackReceived), handlePayPixTransaction(factory))
}

func handlePayPixTransaction(factory *factory.ApplicationFactory) func(context.Context, *event.MessageEvent) error {
	useCase := pay_pix_transaction.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := pay_pix_transaction.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
