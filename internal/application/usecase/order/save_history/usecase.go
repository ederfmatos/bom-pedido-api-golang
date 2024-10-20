package save_history

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/order"
	"context"
	"time"
)

type (
	UseCase struct {
		orderHistoryRepository repository.OrderStatusHistoryRepository
	}

	Input struct {
		Time      time.Time
		OrderId   string
		Status    string
		ChangedBy string
		Data      string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{orderHistoryRepository: factory.OrderHistoryRepository}
}

func (u *UseCase) Execute(ctx context.Context, input Input) error {
	history := &order.StatusHistory{
		Time:      input.Time,
		Status:    input.Status,
		ChangedBy: input.ChangedBy,
		Data:      input.Data,
		OrderId:   input.OrderId,
	}
	return u.orderHistoryRepository.Create(ctx, history)
}
