package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_OrderRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.OrderRepository{
		"OrderMemoryRepository": NewOrderMemoryRepository(),
		"OrderMongoRepository":  NewOrderMongoRepository(container.MongoDatabase()),
	}

	for name, orderRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			customer, err := entity.NewCustomer(faker.Name(), faker.Email(), faker.Word())
			require.NoError(t, err)

			order, err := entity.NewOrder(customer.Id, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 10.0, 100, time.Now(), faker.WORD)
			require.NoError(t, err)

			category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())
			product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, category.Id, faker.WORD)
			require.NoError(t, err)

			err = order.AddProduct(product, 1, "")
			require.NoError(t, err)

			order.Code = 1
			require.NoError(t, err)

			savedOrder, err := orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			require.Nil(t, savedOrder)

			err = orderRepository.Create(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)

			// Approve
			err = order.Approve()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)

			// MarkAsInProgress
			err = order.MarkAsInProgress()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)

			// MarkAsAwaitingDelivery
			err = order.MarkAsAwaitingDelivery()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)

			// MarkAsDelivering
			err = order.MarkAsDelivering()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)

			// Finish
			err = order.Finish()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, order)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, order.Id)
			require.NoError(t, err)
			assertOrder(t, order, savedOrder)
		})
	}
}

func assertOrder(t *testing.T, expectedOrder, actualOrder *entity.Order) {
	require.Equal(t, expectedOrder.Id, actualOrder.Id)
	require.Equal(t, expectedOrder.CustomerID, actualOrder.CustomerID)
	require.Equal(t, expectedOrder.PaymentMethod, actualOrder.PaymentMethod)
	require.Equal(t, expectedOrder.PaymentMode, actualOrder.PaymentMode)
	require.Equal(t, expectedOrder.DeliveryMode, actualOrder.DeliveryMode)
	require.Equal(t, expectedOrder.CreditCardToken, actualOrder.CreditCardToken)
	require.Equal(t, expectedOrder.Payback, actualOrder.Payback)
	require.Equal(t, expectedOrder.Code, actualOrder.Code)
	require.Equal(t, expectedOrder.DeliveryTime.Format("2006-01-02 15:04:05"), actualOrder.DeliveryTime.Format("2006-01-02 15:04:05"))
	require.Equal(t, expectedOrder.GetStatus(), actualOrder.GetStatus())
	require.Equal(t, expectedOrder.Items, actualOrder.Items)
}
