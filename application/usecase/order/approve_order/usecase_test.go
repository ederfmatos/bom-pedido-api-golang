package approve_order

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/order/status"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_ApproveOrder(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := Input{
			OrderId:    value_object.NewID(),
			ApprovedBy: value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should not allow approve order", func(t *testing.T) {
		invalidStatus := []status.Status{
			status.ApprovedStatus,
			status.InProgressStatus,
			status.RejectedStatus,
			status.CancelledStatus,
			status.DeliveringStatus,
			status.AwaitingWithdrawStatus,
			status.AwaitingDeliveryStatus,
			status.FinishedStatus,
		}

		for _, item := range invalidStatus {
			currentStatus := item.Name()
			t.Run("should not allow approve order if status is "+currentStatus, func(t *testing.T) {
				ctx := context.Background()
				orderId := value_object.NewID()
				customerId := value_object.NewID()
				anOrder, err := order.Restore(orderId, customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", currentStatus, time.Now(), 0, 1, 1, time.Now(), []order.Item{}, faker.WORD)
				err = applicationFactory.OrderRepository.Create(ctx, anOrder)
				require.NoError(t, err)
				input := Input{
					OrderId:    anOrder.Id,
					ApprovedBy: value_object.NewID(),
				}
				err = useCase.Execute(ctx, input)
				require.ErrorIs(t, err, status.OperationNotAllowedError)
				savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, anOrder.Id)
				require.NoError(t, err)
				require.Equal(t, savedOrder.GetStatus(), currentStatus)
			})
		}
	})

	t.Run("should approve order", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		anOrder, err := order.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 10.0, time.Now(), faker.WORD)
		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)
		input := Input{
			OrderId:    anOrder.Id,
			ApprovedBy: value_object.NewID(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, anOrder.Id)
		require.NoError(t, err)
		require.Equal(t, savedOrder.GetStatus(), status.ApprovedStatus.Name())
	})
}
