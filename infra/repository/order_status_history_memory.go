package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"context"
)

type OrderStatusHistoryMemoryRepository struct {
	orders map[string][]*order.StatusHistory
}

func NewOrderStatusHistoryMemoryRepository() repository.OrderStatusHistoryRepository {
	return &OrderStatusHistoryMemoryRepository{orders: make(map[string][]*order.StatusHistory)}
}

func (repository *OrderStatusHistoryMemoryRepository) Create(_ context.Context, history *order.StatusHistory) error {
	repository.orders[history.OrderId] = append(repository.orders[history.OrderId], history)
	return nil
}
