package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/internal/domain/entity/order"
	"bom-pedido-api/internal/domain/entity/product"
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

			aCustomer, err := customer.New(faker.Name(), faker.Email(), faker.Word())
			require.NoError(t, err)

			anOrder, err := order.New(aCustomer.Id, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 10.0, 100, time.Now(), faker.WORD)
			require.NoError(t, err)

			category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
			aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.WORD)
			require.NoError(t, err)

			err = anOrder.AddProduct(aProduct, 1, "")
			require.NoError(t, err)

			anOrder.Code = 1
			require.NoError(t, err)

			savedOrder, err := orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			require.Nil(t, savedOrder)

			err = orderRepository.Create(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)

			// Approve
			err = anOrder.Approve()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)

			// MarkAsInProgress
			err = anOrder.MarkAsInProgress()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)

			// MarkAsAwaitingDelivery
			err = anOrder.MarkAsAwaitingDelivery()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)

			// MarkAsDelivering
			err = anOrder.MarkAsDelivering()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)

			// Finish
			err = anOrder.Finish()
			require.NoError(t, err)

			err = orderRepository.Update(ctx, anOrder)
			require.NoError(t, err)

			savedOrder, err = orderRepository.FindById(ctx, anOrder.Id)
			require.NoError(t, err)
			assertOrder(t, anOrder, savedOrder)
		})
	}
}

func assertOrder(t *testing.T, expectedOrder, actualOrder *order.Order) {
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
