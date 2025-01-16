package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/transaction"
	"context"
)

func HandleTransactionEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("CREATE_PIX_TRANSACTION", event.PixPaymentCreated), handleCreatePixTransaction(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("REFUND_PIX_TRANSACTION", event.PixPaymentRefunded), handleRefundPixTransaction(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("CANCEL_PIX_TRANSACTION", event.PixPaymentCancelled), handleCancelPixTransaction(factory))
}

func handleCreatePixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := transaction.NewCreatePixTransaction(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := transaction.CreatePixTransactionInput{
			OrderId:        message.Event.Data["orderId"],
			PaymentId:      message.Event.Data["paymentId"],
			PaymentGateway: message.Event.Data["paymentGateway"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleRefundPixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := transaction.NewRefundPixTransaction(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := transaction.RefundPixTransactionInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleCancelPixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := transaction.NewCancelPixTransaction(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := transaction.CancelPixTransactionInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
