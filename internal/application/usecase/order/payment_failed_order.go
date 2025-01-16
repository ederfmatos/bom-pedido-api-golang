package order

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"context"
	"time"
)

type (
	FailOrderPaymentUseCase struct {
		orderRepository       repository.OrderRepository
		eventEmitter          event.Emitter
		transactionRepository repository.TransactionRepository
	}
	FailOrderPaymentInput struct {
		OrderId string
	}
)

func NewFailOrderPayment(factory *factory.ApplicationFactory) *FailOrderPaymentUseCase {
	return &FailOrderPaymentUseCase{
		orderRepository:       factory.OrderRepository,
		eventEmitter:          factory.EventEmitter,
		transactionRepository: factory.TransactionRepository,
	}
}

func (useCase *FailOrderPaymentUseCase) Execute(ctx context.Context, input FailOrderPaymentInput) error {
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
