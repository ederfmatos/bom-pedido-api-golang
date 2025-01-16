package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type OrderStatusHistoryRepository interface {
	Create(ctx context.Context, history *entity.OrderStatusHistory) error
	ListByOrderId(ctx context.Context, id string) ([]entity.OrderStatusHistory, error)
}
