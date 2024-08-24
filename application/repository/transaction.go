package repository

import (
	"bom-pedido-api/domain/entity/transaction"
	"context"
)

type TransactionRepository interface {
	CreatePixTransaction(ctx context.Context, transaction *transaction.PixTransaction) error
	ExistsByOrderId(ctx context.Context, id string) (bool, error)
}
