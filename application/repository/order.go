package repository

import (
	"bom-pedido-api/domain/entity/order"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *order.Order) error
	FindById(ctx context.Context, id string) (*order.Order, error)
}
