package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/transaction/create_pix_transaction"
	"bom-pedido-api/domain/enums"
	"context"
)

func HandleOrderEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopic("ORDER_CREATED", "CREATE_PIX_TRANSACTION"), handleCreatePixTransaction(factory))
}

func handleCreatePixTransaction(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := create_pix_transaction.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		if message.Event.Data["paymentMethod"] != enums.Pix {
			return message.Ack(ctx)
		}
		input := create_pix_transaction.Input{OrderId: message.Event.Data["orderId"]}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}