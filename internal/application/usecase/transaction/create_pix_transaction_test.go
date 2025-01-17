package transaction

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/gateway/pix"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/mock"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_CreatePixTransaction(t *testing.T) {
	pixGateway := pix.NewFakePixGateway()
	eventEmitter := event.NewMockEventHandler()
	applicationFactory := factory.NewTestApplicationFactory()
	applicationFactory.PixGateway = pixGateway
	applicationFactory.EventEmitter = eventEmitter
	ctx := context.Background()

	customer, err := entity.NewCustomer(faker.Name(), faker.Email(), value_object.NewTenantId())
	require.NoError(t, err)

	customerId := customer.Id
	useCase := NewCreatePixTransaction(applicationFactory)

	t.Run("should return nil if order does not exists", func(t *testing.T) {
		input := CreatePixTransactionInput{OrderId: value_object.NewID()}
		err := useCase.Execute(context.Background(), input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	for _, paymentMethod := range enums.AllPaymentMethods {
		for _, paymentMode := range enums.AllPaymentModes {
			if paymentMethod.IsPix() && paymentMode.IsInApp() {
				continue
			}
			t.Run(fmt.Sprintf("should return nil order is %s %s", paymentMethod.String(), paymentMode.String()), func(t *testing.T) {
				order, err := entity.NewOrder(customerId, paymentMethod.String(), paymentMode.String(), enums.Withdraw, faker.Word(), 0, 0, time.Now(), faker.Word())
				require.NoError(t, err)

				input := CreatePixTransactionInput{OrderId: order.Id, PaymentId: value_object.NewID()}
				err = useCase.Execute(ctx, input)
				require.NoError(t, err)

				eventEmitter.AssertNotCalled(t, "Emit")
			})
		}
	}

	t.Run("should return nil if already exists a transaction to the order", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.TenantId)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		input := CreatePixTransactionInput{OrderId: order.Id, PaymentId: value_object.NewID()}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil payment does not exists", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(nil, nil).Once()

		input := CreatePixTransactionInput{OrderId: order.Id, PaymentId: value_object.NewID()}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{PaymentId: input.PaymentId, MerchantId: order.MerchantId})
	})

	t.Run("should create a pix transaction", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		gatewayPayment := &gateway.GetPaymentOutput{
			Id:             value_object.NewID(),
			QrCode:         faker.Word(),
			ExpiresAt:      time.Now(),
			PaymentGateway: "FAKE",
			QrCodeLink:     faker.URL(),
			Status:         gateway.TransactionPending,
		}
		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(gatewayPayment, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		input := CreatePixTransactionInput{OrderId: order.Id, PaymentId: gatewayPayment.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNumberOfCalls(t, "Emit", 1)
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{PaymentId: input.PaymentId, MerchantId: order.MerchantId})

		savedTransaction, err := applicationFactory.TransactionRepository.FindByOrderId(ctx, order.Id)
		require.NoError(t, err)
		require.Equal(t, savedTransaction.QrCode, gatewayPayment.QrCode)
		require.Equal(t, savedTransaction.QrCodeLink, gatewayPayment.QrCodeLink)
		require.Equal(t, savedTransaction.PaymentGateway, gatewayPayment.PaymentGateway)
		require.Equal(t, savedTransaction.OrderId, order.Id)
		require.Equal(t, savedTransaction.Amount, order.Amount)
	})
}
