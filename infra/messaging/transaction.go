package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/transaction/create_pix_transaction"
	"bom-pedido-api/application/usecase/transaction/refund_pix_transaction"
	"context"
)

func HandleTransactionEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("CREATE_PIX_TRANSACTION", event.PixPaymentCreated), handleCreatePixTransaction(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("REFUND_PIX_TRANSACTION", event.PixPaymentRefunded), handleRefundPixTransaction(factory))
}

func handleCreatePixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := create_pix_transaction.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := create_pix_transaction.Input{
			OrderId:   message.Event.Data["orderId"],
			PaymentId: message.Event.Data["paymentId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleRefundPixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := refund_pix_transaction.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := refund_pix_transaction.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
