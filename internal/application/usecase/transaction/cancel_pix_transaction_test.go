package transaction

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_CancelPixTransaction(t *testing.T) {
	eventEmitter := event.NewMockEventHandler()
	applicationFactory := factory.NewTestApplicationFactory()
	applicationFactory.EventEmitter = eventEmitter
	ctx := context.Background()

	customer, err := entity.NewCustomer(faker.Name(), faker.Email(), value_object.NewTenantId())
	require.NoError(t, err)

	customerId := customer.Id
	useCase := NewCancelPixTransaction(applicationFactory)

	t.Run("should return nil if not exists transaction to the order", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.TenantId)
		require.NoError(t, err)

		err = order.AwaitApproval()
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		input := CancelPixTransactionInput{OrderId: order.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil transaction status is paid", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		pixTransaction.Pay()
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		input := CancelPixTransactionInput{OrderId: order.Id}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil transaction status is cancelled", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		pixTransaction.Cancel()
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		input := CancelPixTransactionInput{OrderId: order.Id}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil transaction status is refunded", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		pixTransaction.Refund()
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		input := CancelPixTransactionInput{OrderId: order.Id}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should refund a pix transaction", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = order.AwaitApproval()
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		input := CancelPixTransactionInput{OrderId: order.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertCalled(t, "Emit", mock.Anything, mock.Anything)

		savedTransaction, err := applicationFactory.TransactionRepository.FindByOrderId(ctx, order.Id)
		require.NoError(t, err)
		require.True(t, savedTransaction.IsCancelled())
	})
}
