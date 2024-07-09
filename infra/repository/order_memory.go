package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"context"
)

type OrderMemoryRepository struct {
	orders map[string]*order.Order
}

func NewOrderMemoryRepository() repository.OrderRepository {
	return &OrderMemoryRepository{orders: make(map[string]*order.Order)}
}

func (repository *OrderMemoryRepository) Create(_ context.Context, order *order.Order) error {
	order.Code = int32(len(repository.orders)) + 1
	repository.orders[order.Id] = order
	return nil
}

func (repository *OrderMemoryRepository) FindById(_ context.Context, id string) (*order.Order, error) {
	return repository.orders[id], nil
}

func (repository *OrderMemoryRepository) Update(_ context.Context, order *order.Order) error {
	repository.orders[order.Id] = order
	return nil
}
