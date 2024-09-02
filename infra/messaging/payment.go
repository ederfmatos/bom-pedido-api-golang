package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/payment/check_pix_payment_failed"
	"bom-pedido-api/application/usecase/payment/create_pix_payment"
	"bom-pedido-api/application/usecase/payment/refund_pix_payment"
	"context"
)

func HandlePaymentEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("CREATE_PIX_PAYMENT", event.OrderCreated), handleCreatePixPayment(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("CHECK_PIX_PAYMENT_FAILED", event.CheckPixPaymentFailed, event.PaymentCallbackReceived), handleCheckPixPaymentFailed(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("REFUND_PIX_PAYMENT", event.OrderCancelled, event.OrderRejected), handleRefundPixPayment(factory))
}

func handleCheckPixPaymentFailed(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := check_pix_payment_failed.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := check_pix_payment_failed.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleCreatePixPayment(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := create_pix_payment.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := create_pix_payment.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleRefundPixPayment(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := refund_pix_payment.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := refund_pix_payment.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
