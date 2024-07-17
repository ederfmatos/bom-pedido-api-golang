package order

import (
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Order(t *testing.T) {
	t.Run("should not allow mark an order as awaiting delivery of order delivery mode is withdraw", func(t *testing.T) {
		customerId := value_object.NewID()
		order, err := New(customerId, enums.CreditCard, enums.InReceiving, enums.Withdraw, "", 0, time.Now())
		assert.NoError(t, err)

		err = order.Approve(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsInProgress(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsAwaitingDelivery(time.Now(), "")
		assert.Error(t, err, errors.OrderDeliveryModeIsWithdrawError)
		err = order.MarkAsAwaitingWithdraw(time.Now(), "")
		assert.NoError(t, err)
	})

	t.Run("should not allow mark an order as awaiting withdraw of order delivery mode is delivery", func(t *testing.T) {
		customerId := value_object.NewID()
		order, err := New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, time.Now())
		assert.NoError(t, err)

		err = order.Approve(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsInProgress(time.Now(), "")
		assert.NoError(t, err)
		err = order.MarkAsAwaitingWithdraw(time.Now(), "")
		assert.Error(t, err, errors.OrderDeliveryModeIsDeliveryError)
		err = order.MarkAsAwaitingDelivery(time.Now(), "")
		assert.NoError(t, err)
	})
}
