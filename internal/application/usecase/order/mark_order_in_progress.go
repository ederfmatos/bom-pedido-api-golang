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
	MarkOrderInProgressUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	MarkOrderInProgressInput struct {
		OrderId string
		By      string
	}
)

func NewMarkOrderInProgress(factory *factory.ApplicationFactory) *MarkOrderInProgressUseCase {
	return &MarkOrderInProgressUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *MarkOrderInProgressUseCase) Execute(ctx context.Context, input MarkOrderInProgressInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.MarkAsInProgress(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderInProgressEvent(order, input.By, time.Now()))
}
