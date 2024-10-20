package await_approval_order

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
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
	if err != nil || order == nil || order.IsAwaitingApproval() {
		return err
	}
	aTransaction, err := useCase.transactionRepository.FindByOrderId(ctx, input.OrderId)
	if err != nil || aTransaction == nil || !aTransaction.IsPaid() {
		return err
	}
	if err = order.AwaitApproval(); err != nil {
		return err
	}
	if err = useCase.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderAwaitingApprovalEvent(order, time.Now()))
}
