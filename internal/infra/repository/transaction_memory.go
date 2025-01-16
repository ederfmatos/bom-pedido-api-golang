package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type TransactionMemoryRepository struct {
	pixTransactionsByOrder map[string]*entity.PixTransaction
}

func NewTransactionMemoryRepository() *TransactionMemoryRepository {
	return &TransactionMemoryRepository{
		pixTransactionsByOrder: make(map[string]*entity.PixTransaction),
	}
}

func (r *TransactionMemoryRepository) CreatePixTransaction(_ context.Context, transaction *entity.PixTransaction) error {
	r.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (r *TransactionMemoryRepository) ExistsByOrderId(_ context.Context, id string) (bool, error) {
	_, ok := r.pixTransactionsByOrder[id]
	return ok, nil
}

func (r *TransactionMemoryRepository) UpdatePixTransaction(_ context.Context, transaction *entity.PixTransaction) error {
	r.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (r *TransactionMemoryRepository) FindByOrderId(_ context.Context, id string) (*entity.PixTransaction, error) {
	return r.pixTransactionsByOrder[id], nil
}
