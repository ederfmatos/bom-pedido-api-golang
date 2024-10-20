package order

import (
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Order(t *testing.T) {
	t.Run("should not allow mark an order as awaiting delivery of order delivery mode is withdraw", func(t *testing.T) {
		customerId := value_object.NewID()
		order, err := New(customerId, enums.CreditCard, enums.InReceiving, enums.Withdraw, "", 0, 0, time.Now(), faker.WORD)
		require.NoError(t, err)

		err = order.Approve()
		require.NoError(t, err)

		err = order.MarkAsInProgress()
		require.NoError(t, err)

		err = order.MarkAsAwaitingDelivery()
		require.Error(t, err, errors.OrderDeliveryModeIsWithdrawError)
		err = order.MarkAsAwaitingWithdraw()
		require.NoError(t, err)

	})

	t.Run("should not allow mark an order as awaiting withdraw of order delivery mode is delivery", func(t *testing.T) {
		customerId := value_object.NewID()
		order, err := New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 0, time.Now(), faker.WORD)
		require.NoError(t, err)

		err = order.Approve()
		require.NoError(t, err)

		err = order.MarkAsInProgress()
		require.NoError(t, err)

		err = order.MarkAsAwaitingWithdraw()
		require.Error(t, err, errors.OrderDeliveryModeIsDeliveryError)
		err = order.MarkAsAwaitingDelivery()
		require.NoError(t, err)
	})
}