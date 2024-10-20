package repository

import (
	"bom-pedido-api/internal/domain/entity/order"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *order.Order) error
	FindById(ctx context.Context, id string) (*order.Order, error)
	Update(ctx context.Context, order *order.Order) error
}
