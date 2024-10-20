package pay_pix_transaction

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/internal/domain/entity/merchant"
	"bom-pedido-api/internal/domain/entity/order"
	"bom-pedido-api/internal/domain/entity/order/status"
	"bom-pedido-api/internal/domain/entity/transaction"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/gateway/pix"
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_PayPixTransaction(t *testing.T) {
	pixGateway := pix.NewFakePixGateway()
	eventEmitter := event.NewMockEventHandler()
	applicationFactory := factory.NewTestApplicationFactory()
	applicationFactory.PixGateway = pixGateway
	applicationFactory.EventEmitter = eventEmitter
	ctx := context.Background()

	aCustomer, err := customer.New(faker.Name(), faker.Email(), value_object.NewTenantId())
	require.NoError(t, err)

	customerId := aCustomer.Id
	useCase := New(applicationFactory)

	t.Run("should return nil if order does not exists", func(t *testing.T) {
		input := Input{OrderId: value_object.NewID()}
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
				anOrder, err := order.New(customerId, paymentMethod.String(), paymentMode.String(), enums.Withdraw, faker.Word(), 0, 0, time.Now(), faker.WORD)
				require.NoError(t, err)

				input := Input{OrderId: anOrder.Id}
				err = useCase.Execute(ctx, input)
				require.NoError(t, err)

				eventEmitter.AssertNotCalled(t, "Emit")
			})
		}
	}

	for _, orderStatus := range status.AllStatus {
		if orderStatus == status.AwaitingPaymentStatus {
			continue
		}
		t.Run("should return nil order is %s"+orderStatus.Name(), func(t *testing.T) {
			anOrder, err := order.Restore(value_object.NewID(), customerId, enums.Pix, enums.InApp, enums.Delivery, "", orderStatus.Name(), time.Now(), 0, 1, 1, time.Now(), []order.Item{}, faker.WORD)
			require.NoError(t, err)

			input := Input{OrderId: anOrder.Id}
			err = useCase.Execute(ctx, input)
			require.NoError(t, err)

			eventEmitter.AssertNotCalled(t, "Emit")
		})
	}

	t.Run("should return nil if not exists transaction to the order", func(t *testing.T) {
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.TenantId)
		require.NoError(t, err)

		err = anOrder.AwaitApproval()
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		input := Input{OrderId: anOrder.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil payment does not exists", func(t *testing.T) {
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		pixTransaction := transaction.NewPixTransaction(value_object.NewID(), anOrder.Id, "", faker.Word(), "", 10)
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(nil, nil).Once()

		input := Input{OrderId: anOrder.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
			PaymentId:      pixTransaction.PaymentId,
			MerchantId:     anOrder.MerchantId,
			PaymentGateway: pixTransaction.PaymentGateway,
		})
	})

	for _, paymentStatus := range []gateway.PaymentStatus{gateway.TransactionCancelled, gateway.TransactionPending, gateway.TransactionRefunded} {
		t.Run("should return nil payment status is "+string(paymentStatus), func(t *testing.T) {
			aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
			require.NoError(t, err)

			err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
			require.NoError(t, err)

			anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.Id)
			require.NoError(t, err)

			err = applicationFactory.OrderRepository.Create(ctx, anOrder)
			require.NoError(t, err)

			pixTransaction := transaction.NewPixTransaction(value_object.NewID(), anOrder.Id, "", faker.Word(), "", 10)
			err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
			require.NoError(t, err)

			pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(&gateway.GetPaymentOutput{
				Id:             value_object.NewID(),
				QrCode:         faker.Word(),
				ExpiresAt:      time.Now(),
				PaymentGateway: "FAKE",
				QrCodeLink:     faker.URL(),
				Status:         paymentStatus,
			}, nil).Once()

			input := Input{OrderId: anOrder.Id}

			err = useCase.Execute(ctx, input)
			require.NoError(t, err)

			eventEmitter.AssertNotCalled(t, "Emit")
			pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
				PaymentId:      pixTransaction.PaymentId,
				MerchantId:     anOrder.MerchantId,
				PaymentGateway: pixTransaction.PaymentGateway,
			})
		})
	}

	t.Run("should pay a pix transaction", func(t *testing.T) {
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.Id)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		pixTransaction := transaction.NewPixTransaction(value_object.NewID(), anOrder.Id, "", faker.Word(), "", 10)
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		gatewayPayment := &gateway.GetPaymentOutput{
			Id:             value_object.NewID(),
			QrCode:         faker.Word(),
			ExpiresAt:      time.Now(),
			PaymentGateway: "FAKE",
			QrCodeLink:     faker.URL(),
			Status:         gateway.TransactionPaid,
		}
		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(gatewayPayment, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		input := Input{OrderId: anOrder.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertCalled(t, "Emit", mock.Anything, mock.Anything)
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
			PaymentId:      pixTransaction.PaymentId,
			MerchantId:     anOrder.MerchantId,
			PaymentGateway: pixTransaction.PaymentGateway,
		})

		savedTransaction, err := applicationFactory.TransactionRepository.FindByOrderId(ctx, anOrder.Id)
		require.NoError(t, err)
		require.True(t, savedTransaction.IsPaid())
	})
}
