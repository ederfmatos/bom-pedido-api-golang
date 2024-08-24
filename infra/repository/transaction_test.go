package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/transaction"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TransactionSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	transactionSqlRepository := NewDefaultTransactionRepository(sqlConnection)
	runTransactionTests(t, transactionSqlRepository)
}

func Test_TransactionMemoryRepository(t *testing.T) {
	transactionSqlRepository := NewTransactionMemoryRepository()
	runTransactionTests(t, transactionSqlRepository)
}

func runTransactionTests(t *testing.T, repository repository.TransactionRepository) {
	ctx := context.TODO()

	aTransaction := transaction.NewPixTransaction(value_object.NewID(), value_object.NewID(), faker.Word(), faker.Word(), faker.Word(), 10)

	existsByOrderId, err := repository.ExistsByOrderId(ctx, aTransaction.OrderId)
	assert.NoError(t, err)
	assert.False(t, existsByOrderId)

	err = repository.CreatePixTransaction(ctx, aTransaction)
	assert.NoError(t, err)

	existsByOrderId, err = repository.ExistsByOrderId(ctx, aTransaction.OrderId)
	assert.NoError(t, err)
	assert.True(t, existsByOrderId)
}
