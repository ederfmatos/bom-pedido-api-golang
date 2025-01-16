package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type OrderStatusHistoryMemoryRepository struct {
	orders map[string][]*entity.OrderStatusHistory
}

func NewOrderStatusHistoryMemoryRepository() *OrderStatusHistoryMemoryRepository {
	return &OrderStatusHistoryMemoryRepository{orders: make(map[string][]*entity.OrderStatusHistory)}
}

func (r *OrderStatusHistoryMemoryRepository) Create(_ context.Context, history *entity.OrderStatusHistory) error {
	r.orders[history.OrderId] = append(r.orders[history.OrderId], history)
	return nil
}

func (r *OrderStatusHistoryMemoryRepository) ListByOrderId(_ context.Context, id string) ([]entity.OrderStatusHistory, error) {
	items := make([]entity.OrderStatusHistory, 0)
	for _, item := range r.orders[id] {
		items = append(items, *item)
	}
	return items, nil
}
