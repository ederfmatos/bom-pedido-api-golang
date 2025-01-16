package order

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/order"
	"context"
	"time"
)

type (
	SaveOrderHistoryUseCase struct {
		orderHistoryRepository repository.OrderStatusHistoryRepository
	}

	SaveOrderHistoryInput struct {
		Time      time.Time
		OrderId   string
		Status    string
		ChangedBy string
		Data      string
	}
)

func NewSaveOrderHistory(factory *factory.ApplicationFactory) *SaveOrderHistoryUseCase {
	return &SaveOrderHistoryUseCase{orderHistoryRepository: factory.OrderHistoryRepository}
}

func (u *SaveOrderHistoryUseCase) Execute(ctx context.Context, input SaveOrderHistoryInput) error {
	history := &order.StatusHistory{
		Time:      input.Time,
		Status:    input.Status,
		ChangedBy: input.ChangedBy,
		Data:      input.Data,
		OrderId:   input.OrderId,
	}
	return u.orderHistoryRepository.Create(ctx, history)
}
