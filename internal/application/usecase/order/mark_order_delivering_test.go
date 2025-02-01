package order

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/entity/status"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
	"time"
)

func Test_MarkOrderDelivering(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewMarkOrderDelivering(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := MarkOrderDeliveringInput{
			OrderId: value_object.NewID(),
			By:      value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should mark an order in delivering", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		order, err := entity.NewOrder(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 0, time.Now(), faker.Word())
		require.NoError(t, err)
		err = order.Approve()
		require.NoError(t, err)

		err = order.MarkAsInProgress()
		require.NoError(t, err)

		err = order.MarkAsAwaitingDelivery()
		require.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)
		input := MarkOrderDeliveringInput{
			OrderId: order.Id,
			By:      value_object.NewID(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
		require.NoError(t, err)
		require.Equal(t, savedOrder.GetStatus(), status.DeliveringStatus.Name())
	})
}
