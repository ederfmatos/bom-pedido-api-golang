package order

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/errors"
	"context"
	"time"
)

type (
	ApproveOrderUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	ApproveOrderUseCaseInput struct {
		OrderId    string
		ApprovedBy string
	}
)

func NewApproveOrderUseCase(factory *factory.ApplicationFactory) *ApproveOrderUseCase {
	return &ApproveOrderUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *ApproveOrderUseCase) Execute(ctx context.Context, input ApproveOrderUseCaseInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.Approve(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderApprovedEvent(order, input.ApprovedBy, time.Now()))
}
