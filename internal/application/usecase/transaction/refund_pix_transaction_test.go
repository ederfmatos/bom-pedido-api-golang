package transaction

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/entity/status"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/gateway/pix"
	"bom-pedido-api/internal/infra/lock"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/mock"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_RefundPixTransaction(t *testing.T) {
	var (
		ctx = context.Background()

		pixGateway            = pix.NewFakePixGateway()
		eventEmitter          = event.NewMockEventHandler()
		merchantRepository    = repository.NewMerchantMemoryRepository()
		orderRepository       = repository.NewOrderMemoryRepository()
		transactionRepository = repository.NewTransactionMemoryRepository()
		locker                = lock.NewMemoryLocker()
	)

	customer, err := entity.NewCustomer(faker.Name(), faker.Email(), value_object.NewTenantId())
	require.NoError(t, err)

	customerId := customer.Id
	useCase := NewRefundPixTransaction(
		orderRepository,
		transactionRepository,
		pixGateway,
		eventEmitter,
		locker,
	)

	t.Run("should return nil if order does not exists", func(t *testing.T) {
		input := RefundPixTransactionInput{OrderId: value_object.NewID()}
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

				input := RefundPixTransactionInput{OrderId: order.Id}
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
			order, err := entity.RestoreOrder(value_object.NewID(), customerId, enums.Pix, enums.InApp, enums.Delivery, "", orderStatus.Name(), time.Now(), 0, 1, 1, time.Now(), []entity.OrderItem{}, faker.Word())
			require.NoError(t, err)

			input := RefundPixTransactionInput{OrderId: order.Id}
			err = useCase.Execute(ctx, input)
			require.NoError(t, err)

			eventEmitter.AssertNotCalled(t, "Emit")
		})
	}

	t.Run("should return nil if not exists transaction to the order", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = merchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.TenantId)
		require.NoError(t, err)

		err = order.AwaitApproval()
		require.NoError(t, err)

		err = orderRepository.Create(ctx, order)
		require.NoError(t, err)

		input := RefundPixTransactionInput{OrderId: order.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
	})

	t.Run("should return nil payment does not exists", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = merchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = order.AwaitApproval()
		require.NoError(t, err)

		err = orderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		pixTransaction.Pay()
		err = transactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(nil, nil).Once()

		input := RefundPixTransactionInput{OrderId: order.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNotCalled(t, "Emit")
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
			PaymentId:      pixTransaction.PaymentId,
			MerchantId:     order.MerchantId,
			PaymentGateway: pixTransaction.PaymentGateway,
		})
	})

	for _, transactionStatus := range []string{"CREATED", "REFUNDED"} {
		t.Run("should return nil pix transaction is %s"+transactionStatus, func(t *testing.T) {
			merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
			require.NoError(t, err)

			err = merchantRepository.Create(ctx, merchant)
			require.NoError(t, err)

			order, err := entity.RestoreOrder(value_object.NewID(), customerId, enums.Pix, enums.InApp, enums.Delivery, "", status.AwaitingApprovalStatus.Name(), time.Now(), 0, 1, 1, time.Now(), []entity.OrderItem{}, faker.Word())
			require.NoError(t, err)

			err = orderRepository.Create(ctx, order)
			require.NoError(t, err)

			pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
			pixTransaction.Status = entity.PixTransactionStatus(transactionStatus)
			err = transactionRepository.CreatePixTransaction(ctx, pixTransaction)
			require.NoError(t, err)

			input := RefundPixTransactionInput{OrderId: order.Id}
			err = useCase.Execute(ctx, input)
			require.NoError(t, err)

			eventEmitter.AssertNotCalled(t, "Emit")
			pixGateway.AssertNotCalled(t, "GetPaymentById")
		})
	}

	for _, paymentStatus := range []gateway.PaymentStatus{gateway.TransactionCancelled, gateway.TransactionPending, gateway.TransactionPaid} {
		t.Run("should return nil payment status is "+string(paymentStatus), func(t *testing.T) {
			merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
			require.NoError(t, err)

			err = merchantRepository.Create(ctx, merchant)
			require.NoError(t, err)

			order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
			require.NoError(t, err)

			err = order.AwaitApproval()
			require.NoError(t, err)

			err = orderRepository.Create(ctx, order)
			require.NoError(t, err)

			pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
			pixTransaction.Pay()
			err = transactionRepository.CreatePixTransaction(ctx, pixTransaction)
			require.NoError(t, err)

			pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(&gateway.GetPaymentOutput{
				Id:             value_object.NewID(),
				QrCode:         faker.Word(),
				ExpiresAt:      time.Now(),
				PaymentGateway: "FAKE",
				QrCodeLink:     faker.URL(),
				Status:         paymentStatus,
			}, nil).Once()

			input := RefundPixTransactionInput{OrderId: order.Id}

			err = useCase.Execute(ctx, input)
			require.NoError(t, err)

			eventEmitter.AssertNotCalled(t, "Emit")
			pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
				PaymentId:      pixTransaction.PaymentId,
				MerchantId:     order.MerchantId,
				PaymentGateway: pixTransaction.PaymentGateway,
			})
		})
	}

	t.Run("should refund a pix transaction", func(t *testing.T) {
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = merchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		order, err := entity.NewOrder(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), merchant.Id)
		require.NoError(t, err)

		err = order.AwaitApproval()
		require.NoError(t, err)

		err = orderRepository.Create(ctx, order)
		require.NoError(t, err)

		pixTransaction := entity.NewPixTransaction(value_object.NewID(), order.Id, "", faker.Word(), "", 10)
		pixTransaction.Pay()
		err = transactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		gatewayPayment := &gateway.GetPaymentOutput{
			Id:             value_object.NewID(),
			QrCode:         faker.Word(),
			ExpiresAt:      time.Now(),
			PaymentGateway: "FAKE",
			QrCodeLink:     faker.URL(),
			Status:         gateway.TransactionRefunded,
		}
		pixGateway.On("GetPaymentById", mock.Anything, mock.Anything).Return(gatewayPayment, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		input := RefundPixTransactionInput{OrderId: order.Id}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertCalled(t, "Emit", mock.Anything, mock.Anything)
		pixGateway.AssertCalled(t, "GetPaymentById", ctx, gateway.GetPaymentInput{
			PaymentId:      pixTransaction.PaymentId,
			MerchantId:     order.MerchantId,
			PaymentGateway: pixTransaction.PaymentGateway,
		})

		savedTransaction, err := transactionRepository.FindByOrderId(ctx, order.Id)
		require.NoError(t, err)
		require.True(t, savedTransaction.IsRefunded())
	})
}
