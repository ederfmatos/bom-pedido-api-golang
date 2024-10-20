package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/transaction"
	"context"
)

type TransactionMemoryRepository struct {
	pixTransactionsByOrder map[string]*transaction.PixTransaction
}

func NewTransactionMemoryRepository() repository.TransactionRepository {
	return &TransactionMemoryRepository{
		pixTransactionsByOrder: make(map[string]*transaction.PixTransaction),
	}
}

func (repository *TransactionMemoryRepository) CreatePixTransaction(_ context.Context, transaction *transaction.PixTransaction) error {
	repository.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (repository *TransactionMemoryRepository) ExistsByOrderId(_ context.Context, id string) (bool, error) {
	_, ok := repository.pixTransactionsByOrder[id]
	return ok, nil
}

func (repository *TransactionMemoryRepository) UpdatePixTransaction(_ context.Context, transaction *transaction.PixTransaction) error {
	repository.pixTransactionsByOrder[transaction.OrderId] = transaction
	return nil
}

func (repository *TransactionMemoryRepository) FindByOrderId(_ context.Context, id string) (*transaction.PixTransaction, error) {
	return repository.pixTransactionsByOrder[id], nil
}
