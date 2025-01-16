package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type TransactionRepository interface {
	CreatePixTransaction(ctx context.Context, transaction *entity.PixTransaction) error
	UpdatePixTransaction(ctx context.Context, transaction *entity.PixTransaction) error
	ExistsByOrderId(ctx context.Context, id string) (bool, error)
	FindByOrderId(ctx context.Context, id string) (*entity.PixTransaction, error)
}
