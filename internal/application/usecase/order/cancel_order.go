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
	CancelOrderUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	CancelOrderInput struct {
		OrderId     string
		CancelledBy string
		Reason      string
	}
)

func NewCancelOrder(factory *factory.ApplicationFactory) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *CancelOrderUseCase) Execute(ctx context.Context, input CancelOrderInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.Cancel(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderCancelledEvent(order, input.CancelledBy, time.Now(), input.Reason))
}
