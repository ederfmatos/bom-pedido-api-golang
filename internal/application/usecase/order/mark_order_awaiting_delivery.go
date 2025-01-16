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
	MarkOrderAwaitingDeliveryUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	MarkOrderAwaitingDeliveryInput struct {
		OrderId string
		By      string
	}
)

func NewMarkOrderAwaitingDelivery(factory *factory.ApplicationFactory) *MarkOrderAwaitingDeliveryUseCase {
	return &MarkOrderAwaitingDeliveryUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *MarkOrderAwaitingDeliveryUseCase) Execute(ctx context.Context, input MarkOrderAwaitingDeliveryInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.MarkAsAwaitingDelivery(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderAwaitingDeliveryEvent(order, input.By, time.Now()))
}
