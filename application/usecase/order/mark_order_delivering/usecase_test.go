package mark_order_delivering

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
			OrderId: value_object.NewID(),
			By:      value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		assert.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should mark an order in delivering", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		order, err := order2.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 0, time.Now(), faker.WORD)
		err = order.Approve(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsInProgress(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsAwaitingDelivery(time.Now(), "")
		assert.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, order)
		assert.NoError(t, err)
		input := Input{
			OrderId: order.Id,
			By:      value_object.NewID(),
		}
		err = useCase.Execute(ctx, input)
		assert.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, order.Id)
		assert.NoError(t, err)
		assert.Equal(t, savedOrder.GetStatus(), status.DeliveringStatus.Name())
	})
}
