package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TransactionRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.TransactionRepository{
		"TransactionMemoryRepository": NewTransactionMemoryRepository(),
		"TransactionMongoRepository":  NewTransactionMongoRepository(container.MongoDatabase()),
	}

	for name, transactionRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			pixTransaction := entity.NewPixTransaction(value_object.NewID(), value_object.NewID(), faker.Word(), faker.Word(), faker.Word(), 10)

			existsByOrderId, err := transactionRepository.ExistsByOrderId(ctx, pixTransaction.OrderId)
			require.NoError(t, err)
			require.False(t, existsByOrderId)

			err = transactionRepository.CreatePixTransaction(ctx, pixTransaction)
			require.NoError(t, err)

			existsByOrderId, err = transactionRepository.ExistsByOrderId(ctx, pixTransaction.OrderId)
			require.NoError(t, err)
			require.True(t, existsByOrderId)

			savedPixTransaction, err := transactionRepository.FindByOrderId(ctx, pixTransaction.OrderId)
			require.NoError(t, err)
			require.NotNil(t, savedPixTransaction)
			require.Equal(t, pixTransaction, savedPixTransaction)

			pixTransaction.Cancel()

			err = transactionRepository.UpdatePixTransaction(ctx, pixTransaction)
			require.NoError(t, err)

			savedPixTransaction, err = transactionRepository.FindByOrderId(ctx, pixTransaction.OrderId)
			require.NoError(t, err)
			require.NotNil(t, savedPixTransaction)
			require.Equal(t, pixTransaction, savedPixTransaction)
		})
	}
}
