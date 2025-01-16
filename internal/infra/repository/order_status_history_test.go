package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOrderStatusHistoryMongoRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.OrderStatusHistoryRepository{
		"OrderStatusHistoryMongoRepository": NewOrderStatusHistoryMongoRepository(container.MongoDatabase()),
	}

	for name, orderStatusHistoryRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			orderID := "orderID"

			history, err := orderStatusHistoryRepository.ListByOrderId(ctx, orderID)
			require.NoError(t, err)
			require.Empty(t, history)

			statusHistory := entity.NewOrderStatusHistory(time.Now(), "CREATED", faker.Word(), faker.Word(), orderID)

			err = orderStatusHistoryRepository.Create(ctx, statusHistory)
			require.NoError(t, err)

			history, err = orderStatusHistoryRepository.ListByOrderId(ctx, orderID)
			require.NoError(t, err)
			require.NotEmpty(t, history)
			require.Equal(t, *statusHistory, history[0])
		})
	}
}
