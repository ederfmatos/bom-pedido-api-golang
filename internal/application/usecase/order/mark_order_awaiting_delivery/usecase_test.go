package mark_order_awaiting_delivery

import (
	"bom-pedido-api/internal/domain/entity/order"
	"bom-pedido-api/internal/domain/entity/order/status"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_MarkOrderAwaitingDelivery(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := Input{
			OrderId: value_object.NewID(),
			By:      value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should mark an order in delivering", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		anOrder, err := order.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 0, time.Now(), faker.WORD)
		require.NoError(t, err)
		err = anOrder.Approve()
		require.NoError(t, err)

		err = anOrder.MarkAsInProgress()
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)
		input := Input{
			OrderId: anOrder.Id,
			By:      value_object.NewID(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
		savedOrder, err := applicationFactory.OrderRepository.FindById(ctx, anOrder.Id)
		require.NoError(t, err)
		require.Equal(t, savedOrder.GetStatus(), status.AwaitingDeliveryStatus.Name())
	})
}
