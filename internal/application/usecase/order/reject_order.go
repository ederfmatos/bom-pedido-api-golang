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
	RejectOrderUseCase struct {
		orderRepository repository.OrderRepository
		eventEmitter    event.Emitter
	}
	RejectOrderInput struct {
		OrderId    string
		RejectedBy string
		Reason     string
	}
)

func NewRejectOrder(factory *factory.ApplicationFactory) *RejectOrderUseCase {
	return &RejectOrderUseCase{
		orderRepository: factory.OrderRepository,
		eventEmitter:    factory.EventEmitter,
	}
}

func (useCase *RejectOrderUseCase) Execute(ctx context.Context, input RejectOrderInput) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	if err = order.Reject(); err != nil {
		return err
	}
	err = useCase.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return useCase.eventEmitter.Emit(ctx, event.NewOrderRejectedEvent(order, input.RejectedBy, time.Now(), input.Reason))
}
