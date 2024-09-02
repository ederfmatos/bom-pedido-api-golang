package payment_failed_order

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		eventEmitter          event.Emitter
		transactionRepository repository.TransactionRepository
	}
	Input struct {
		OrderId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository:       factory.OrderRepository,
		eventEmitter:          factory.EventEmitter,
		transactionRepository: factory.TransactionRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsAwaitingPayment() {
		return err
	}
	if err = order.PaymentFailed(); err != nil {
		return err
	}
	if err = useCase.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderPaymentFailedEvent(order, time.Now()))
}
