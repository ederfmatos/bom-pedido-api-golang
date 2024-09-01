package mark_order_delivering

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	Input struct {
		OrderId string
		By      string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
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
