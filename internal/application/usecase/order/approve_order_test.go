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

func Test_ApproveOrder(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewApproveOrder(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := ApproveOrderUseCaseInput{
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
				order, err := entity.RestoreOrder(orderId, customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", currentStatus, time.Now(), 0, 1, 1, time.Now(), []entity.OrderItem{}, faker.Word())
				require.NoError(t, err)
				err = applicationFactory.OrderRepository.Create(ctx, order)
				require.NoError(t, err)
				input := ApproveOrderUseCaseInput{
					OrderId:    order.Id,
					ApprovedBy: value_object.NewID(),
				}
				err = useCase.Execute(ctx, input)
				require.ErrorIs(t, err, status.OperationNotAllowedError)
				savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
				require.NoError(t, err)
				require.Equal(t, savedOrder.GetStatus(), currentStatus)
			})
		}
	})

	t.Run("should approve order", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		order, err := entity.NewOrder(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 10.0, time.Now(), faker.Word())
		require.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)
		input := ApproveOrderUseCaseInput{
			OrderId:    order.Id,
			ApprovedBy: value_object.NewID(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
		require.NoError(t, err)
		require.Equal(t, savedOrder.GetStatus(), status.ApprovedStatus.Name())
	})
}
