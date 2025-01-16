package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/payment"
	"context"
)

func HandlePaymentEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("CREATE_PIX_PAYMENT", event.OrderCreated), handleCreatePixPayment(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("CHECK_PIX_PAYMENT_FAILED", event.CheckPixPaymentFailed, event.PaymentCallbackReceived), handleCheckPixPaymentFailed(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("REFUND_PIX_PAYMENT", event.OrderCancelled, event.OrderRejected), handleRefundPixPayment(factory))
}

func handleCheckPixPaymentFailed(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := payment.NewCheckPixPaymentFailed(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := payment.CheckPixPaymentFailedInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleCreatePixPayment(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := payment.NewCreatePixPayment(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := payment.CreatePixPaymentInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleRefundPixPayment(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := payment.NewRefundPixPayment(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := payment.RefundPixPaymentInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
