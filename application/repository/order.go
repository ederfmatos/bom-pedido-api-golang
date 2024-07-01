package repository

import (
	"bom-pedido-api/domain/entity"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
}
