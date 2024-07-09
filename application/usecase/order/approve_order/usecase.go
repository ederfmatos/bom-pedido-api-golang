package approve_order

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/events"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	Input struct {
		Context    context.Context
		OrderId    string
		ApprovedBy string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *UseCase) Execute(input Input) error {
	order, err := useCase.orderRepository.FindById(input.Context, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	err = order.Approve(time.Now(), input.ApprovedBy)
	if err != nil {
		return err
	}
	err = useCase.orderRepository.Update(input.Context, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(input.Context, events.NewOrderApprovedEvent(order))
}
