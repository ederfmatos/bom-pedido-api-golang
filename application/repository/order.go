package repository

import (
	"bom-pedido-api/domain/entity"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	FindById(ctx context.Context, id string) (*entity.Order, error)
}