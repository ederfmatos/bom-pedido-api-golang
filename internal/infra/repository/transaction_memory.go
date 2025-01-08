package repository

import (
	"bom-pedido-api/internal/domain/entity/transaction"
	"context"
)

type TransactionMemoryRepository struct {
	pixTransactionsByOrder map[string]*transaction.PixTransaction
}

func NewTransactionMemoryRepository() *TransactionMemoryRepository {
	return &TransactionMemoryRepository{
		pixTransactionsByOrder: make(map[string]*transaction.PixTransaction),
	}
}

func (r *TransactionMemoryRepository) CreatePixTransaction(_ context.Context, transaction *transaction.PixTransaction) error {
	r.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (r *TransactionMemoryRepository) ExistsByOrderId(_ context.Context, id string) (bool, error) {
	_, ok := r.pixTransactionsByOrder[id]
	return ok, nil
}

func (r *TransactionMemoryRepository) UpdatePixTransaction(_ context.Context, transaction *transaction.PixTransaction) error {
	r.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (r *TransactionMemoryRepository) FindByOrderId(_ context.Context, id string) (*transaction.PixTransaction, error) {
	return r.pixTransactionsByOrder[id], nil
}
