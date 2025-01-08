package repository

import (
	"bom-pedido-api/internal/domain/entity/transaction"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type TransactionMongoRepository struct {
	collection *mongo.Collection
}

func NewTransactionMongoRepository(database *mongo.Database) *TransactionMongoRepository {
	return &TransactionMongoRepository{collection: database.ForCollection("transactions")}
}

func (r *TransactionMongoRepository) CreatePixTransaction(ctx context.Context, transaction *transaction.PixTransaction) error {
	return r.collection.InsertOne(ctx, transaction)
}

func (r *TransactionMongoRepository) ExistsByOrderId(ctx context.Context, id string) (bool, error) {
	return r.collection.ExistsBy(ctx, "orderId", id)
}

func (r *TransactionMongoRepository) UpdatePixTransaction(ctx context.Context, transaction *transaction.PixTransaction) error {
	return r.collection.UpdateByID(ctx, transaction.Id, transaction)
}

func (r *TransactionMongoRepository) FindByOrderId(ctx context.Context, id string) (*transaction.PixTransaction, error) {
	var pixTransaction transaction.PixTransaction
	err := r.collection.FindBy(ctx, "orderId", id, &pixTransaction)
	if err != nil || pixTransaction.Id == "" {
		return nil, err
	}
	return &pixTransaction, nil
}
