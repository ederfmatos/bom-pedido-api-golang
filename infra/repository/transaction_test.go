package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/transaction"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
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

	pixTransaction := transaction.NewPixTransaction(value_object.NewID(), value_object.NewID(), faker.Word(), faker.Word(), faker.Word(), 10)

	existsByOrderId, err := repository.ExistsByOrderId(ctx, pixTransaction.OrderId)
	require.NoError(t, err)
	require.False(t, existsByOrderId)

	err = repository.CreatePixTransaction(ctx, pixTransaction)
	require.NoError(t, err)

	existsByOrderId, err = repository.ExistsByOrderId(ctx, pixTransaction.OrderId)
	require.NoError(t, err)
	require.True(t, existsByOrderId)
}
