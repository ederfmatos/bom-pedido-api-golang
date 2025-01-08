package repository

import (
	"bom-pedido-api/internal/domain/entity/order"
	"context"
)

type OrderStatusHistoryMemoryRepository struct {
	orders map[string][]*order.StatusHistory
}

func NewOrderStatusHistoryMemoryRepository() *OrderStatusHistoryMemoryRepository {
	return &OrderStatusHistoryMemoryRepository{orders: make(map[string][]*order.StatusHistory)}
}

func (r *OrderStatusHistoryMemoryRepository) Create(_ context.Context, history *order.StatusHistory) error {
	r.orders[history.OrderId] = append(r.orders[history.OrderId], history)
	return nil
}

func (r *OrderStatusHistoryMemoryRepository) ListByOrderId(_ context.Context, id string) ([]order.StatusHistory, error) {
	items := make([]order.StatusHistory, 0)
	for _, item := range r.orders[id] {
		items = append(items, *item)
	}
	return items, nil
}
