package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type OrderMemoryRepository struct {
	orders map[string]*entity.Order
}

func NewOrderMemoryRepository() repository.OrderRepository {
	return &OrderMemoryRepository{orders: make(map[string]*entity.Order)}
}

func (repository *OrderMemoryRepository) Create(_ context.Context, order *entity.Order) error {
	order.Code = int32(len(repository.orders)) + 1
	repository.orders[order.Id] = order
	return nil
}

func (repository *OrderMemoryRepository) FindById(_ context.Context, id string) (*entity.Order, error) {
	return repository.orders[id], nil
}
