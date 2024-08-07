package cancel_order

import (
	order2 "bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/order/status"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_UseCase(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := Input{
			OrderId:     value_object.NewID(),
			CancelledBy: value_object.NewID(),
			Reason:      faker.Word(),
		}
		err := useCase.Execute(ctx, input)
		assert.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should cancel order", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		order, err := order2.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, time.Now())
		err = order.Approve(time.Now(), "")
		assert.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, order)
		assert.NoError(t, err)
		input := Input{
			OrderId:     order.Id,
			CancelledBy: value_object.NewID(),
			Reason:      faker.Word(),
		}
		err = useCase.Execute(ctx, input)
		assert.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
		assert.NoError(t, err)
		assert.Equal(t, savedOrder.GetStatus(), status.CancelledStatus.Name())
	})

	t.Run("should not allow approve order", func(t *testing.T) {
		invalidStatus := []status.Status{
			status.AwaitingApprovalStatus,
			status.CancelledStatus,
			status.RejectedStatus,
		}

		for _, item := range invalidStatus {
			currentStatus := item.Name()
			t.Run("should not allow cancel order if status is "+currentStatus, func(t *testing.T) {
				ctx := context.Background()
				orderId := value_object.NewID()
				customerId := value_object.NewID()
				order, err := order2.Restore(orderId, customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", currentStatus, time.Now(), 0, 1, time.Now(), []order2.Item{}, make([]status.History, 0))
				err = applicationFactory.OrderRepository.Create(ctx, order)
				assert.NoError(t, err)
				input := Input{
					OrderId:     order.Id,
					CancelledBy: value_object.NewID(),
					Reason:      faker.Word(),
				}
				err = useCase.Execute(ctx, input)
				assert.ErrorIs(t, err, status.OperationNotAllowedError)
				savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
				assert.NoError(t, err)
				assert.Equal(t, savedOrder.GetStatus(), currentStatus)
			})
		}
	})
}
