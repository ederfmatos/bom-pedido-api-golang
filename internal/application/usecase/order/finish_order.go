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
	FinishOrderUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	FinishOrderInput struct {
		OrderId    string
		FinishedBy string
	}
)

func NewFinishOrder(factory *factory.ApplicationFactory) *FinishOrderUseCase {
	return &FinishOrderUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *FinishOrderUseCase) Execute(ctx context.Context, input FinishOrderInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.Finish(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderFinishedEvent(order, input.FinishedBy, time.Now()))
}
