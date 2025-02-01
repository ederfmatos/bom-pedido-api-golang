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
	MarkOrderDeliveringUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	MarkOrderDeliveringInput struct {
		OrderId string
		By      string
	}
)

func NewMarkOrderDelivering(factory *factory.ApplicationFactory) *MarkOrderDeliveringUseCase {
	return &MarkOrderDeliveringUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *MarkOrderDeliveringUseCase) Execute(ctx context.Context, input MarkOrderDeliveringInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.MarkAsDelivering(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderDeliveringEvent(order, input.By, time.Now()))
}
