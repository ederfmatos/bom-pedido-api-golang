package repository

import (
	"bom-pedido-api/internal/domain/entity/order"
	"context"
)

type OrderMemoryRepository struct {
	orders map[string]*order.Order
}

func NewOrderMemoryRepository() *OrderMemoryRepository {
	return &OrderMemoryRepository{orders: make(map[string]*order.Order)}
}

func (r *OrderMemoryRepository) Create(_ context.Context, order *order.Order) error {
	order.Code = int32(len(r.orders)) + 1
	r.orders[order.Id] = order
	return nil
}

func (r *OrderMemoryRepository) FindById(_ context.Context, id string) (*order.Order, error) {
	return r.orders[id], nil
}

func (r *OrderMemoryRepository) Update(_ context.Context, order *order.Order) error {
	r.orders[order.Id] = order
	return nil
}
