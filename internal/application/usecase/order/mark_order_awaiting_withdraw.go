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
	MarkOrderAwaitingWithdrawUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	MarkOrderAwaitingWithdrawInput struct {
		OrderId string
		By      string
	}
)

func NewMarkOrderAwaitingWithdraw(factory *factory.ApplicationFactory) *MarkOrderAwaitingWithdrawUseCase {
	return &MarkOrderAwaitingWithdrawUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *MarkOrderAwaitingWithdrawUseCase) Execute(ctx context.Context, input MarkOrderAwaitingWithdrawInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.MarkAsAwaitingWithdraw(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderAwaitingWithdrawEvent(order, input.By, time.Now()))
}
