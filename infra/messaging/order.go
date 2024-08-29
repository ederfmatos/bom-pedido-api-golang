package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/await_approval_order"
	"context"
)

func HandleOrderEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("AWAIT_APPROVAL_ORDER", event.PixTransactionPaid), handleAwaitApprovalOrder(factory))
}

func handleAwaitApprovalOrder(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := await_approval_order.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := await_approval_order.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
