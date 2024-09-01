package repository

import (
	"bom-pedido-api/domain/entity/order"
	"context"
)

type OrderStatusHistoryRepository interface {
	Create(ctx context.Context, history *order.StatusHistory) error
}
