package repository

import (
	"bom-pedido-api/internal/domain/entity/order"
	"context"
)

type OrderStatusHistoryRepository interface {
	Create(ctx context.Context, history *order.StatusHistory) error
	ListByOrderId(ctx context.Context, id string) ([]order.StatusHistory, error)
}
