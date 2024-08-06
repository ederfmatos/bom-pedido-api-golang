package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OrderSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	orderRepository := NewDefaultOrderRepository(sqlConnection)
	orderTests(t, orderRepository)
}

func Test_OrderMemoryRepository(t *testing.T) {
	orderRepository := NewOrderMemoryRepository()
	orderTests(t, orderRepository)
}

func orderTests(t *testing.T, repository repository.OrderRepository) {
	ctx := context.TODO()

	customerId := value_object.NewID()
	adminId := value_object.NewID()
	anOrder, err := order.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 10.0, time.Now())
	anOrder.Code = 1
	assert.NoError(t, err)

	savedOrder, err := repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Nil(t, savedOrder)

	err = repository.Create(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)

	// Approve
	err = anOrder.Approve(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)

	// MarkAsInProgress
	err = anOrder.MarkAsInProgress(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)

	// MarkAsAwaitingDelivery
	err = anOrder.MarkAsAwaitingDelivery(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)

	// MarkAsDelivering
	err = anOrder.MarkAsDelivering(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)

	// Finish
	err = anOrder.Finish(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assertOrder(t, anOrder, savedOrder)
}

func assertOrder(t *testing.T, expectedOrder, actualOrder *order.Order) {
	assert.Equal(t, expectedOrder.Id, actualOrder.Id)
	assert.Equal(t, expectedOrder.CustomerID, actualOrder.CustomerID)
	assert.Equal(t, expectedOrder.PaymentMethod, actualOrder.PaymentMethod)
	assert.Equal(t, expectedOrder.PaymentMode, actualOrder.PaymentMode)
	assert.Equal(t, expectedOrder.DeliveryMode, actualOrder.DeliveryMode)
	assert.Equal(t, expectedOrder.CreditCardToken, actualOrder.CreditCardToken)
	assert.Equal(t, expectedOrder.Change, actualOrder.Change)
	assert.Equal(t, expectedOrder.Code, actualOrder.Code)
	assert.Equal(t, expectedOrder.DeliveryTime.Format("2006-01-02 15:04:05"), actualOrder.DeliveryTime.Format("2006-01-02 15:04:05"))
	assert.Equal(t, expectedOrder.GetStatus(), actualOrder.GetStatus())
	assert.Equal(t, expectedOrder.Items, actualOrder.Items)
	assert.Equal(t, expectedOrder.History, actualOrder.History)
}
