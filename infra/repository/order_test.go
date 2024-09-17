package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_OrderSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	orderRepository := NewDefaultOrderRepository(sqlConnection)
	productRepository := NewDefaultProductRepository(sqlConnection)
	customerRepository := NewDefaultCustomerRepository(sqlConnection)
	categoryRepository := NewDefaultProductCategoryRepository(sqlConnection)
	orderTests(t, orderRepository, categoryRepository, productRepository, customerRepository)
}

func Test_OrderMemoryRepository(t *testing.T) {
	orderRepository := NewOrderMemoryRepository()
	productRepository := NewProductMemoryRepository()
	categoryRepository := NewProductCategoryMemoryRepository()
	customerRepository := NewCustomerMemoryRepository()
	orderTests(t, orderRepository, categoryRepository, productRepository, customerRepository)
}

func orderTests(
	t *testing.T,
	orderRepository repository.OrderRepository,
	categoryRepository repository.ProductCategoryRepository,
	productRepository repository.ProductRepository,
	customerRepository repository.CustomerRepository,
) {
	ctx := context.Background()

	aCustomer, err := customer.New(faker.Name(), faker.Email(), faker.Word())
	require.NoError(t, err)

	err = customerRepository.Create(ctx, aCustomer)
	require.NoError(t, err)

	anOrder, err := order.New(aCustomer.Id, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 10.0, 100, time.Now(), faker.WORD)
	require.NoError(t, err)

	category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
	err = categoryRepository.Create(ctx, category)
	require.NoError(t, err)

	aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.WORD)
	require.NoError(t, err)

	err = anOrder.AddProduct(aProduct, 1, "")
	require.NoError(t, err)

	err = productRepository.Create(ctx, aProduct)
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
