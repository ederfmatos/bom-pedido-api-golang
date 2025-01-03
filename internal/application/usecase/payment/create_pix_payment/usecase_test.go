package create_pix_payment

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/internal/domain/entity/merchant"
	"bom-pedido-api/internal/domain/entity/order"
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

func Test_CreatePixPayment(t *testing.T) {
	pixGateway := pix.NewFakePixGateway()
	eventEmitter := event.NewMockEventHandler()
	applicationFactory := factory.NewTestApplicationFactory()
	applicationFactory.PixGateway = pixGateway
	applicationFactory.EventEmitter = eventEmitter

	aCustomer, err := customer.New(faker.Name(), faker.Email(), value_object.NewTenantId())
	require.NoError(t, err)
	err = applicationFactory.CustomerRepository.Create(context.Background(), aCustomer)
	require.NoError(t, err)

	customerId := aCustomer.Id
	t.Run("should return nil if order does not exists", func(t *testing.T) {
		useCase := New(applicationFactory)
		input := Input{OrderId: value_object.NewID()}

		err := useCase.Execute(context.Background(), input)
		require.NoError(t, err)
	})

	for _, paymentMethod := range enums.AllPaymentMethods {
		for _, paymentMode := range enums.AllPaymentModes {
			if paymentMethod.IsPix() && paymentMode.IsInApp() {
				continue
			}
			t.Run(fmt.Sprintf("should return nil order is %s %s", paymentMethod.String(), paymentMode.String()), func(t *testing.T) {
				anOrder, err := order.New(customerId, paymentMethod.String(), paymentMode.String(), enums.Withdraw, faker.Word(), 0, 0, time.Now(), faker.WORD)
				require.NoError(t, err)

				useCase := New(applicationFactory)
				input := Input{OrderId: anOrder.Id}

				err = useCase.Execute(context.Background(), input)
				require.NoError(t, err)
			})
		}
	}

	t.Run("should return nil if customer does not exists", func(t *testing.T) {
		ctx := context.Background()

		anOrder, err := order.New(value_object.NewID(), enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 0, time.Now(), faker.WORD)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		useCase := New(applicationFactory)
		input := Input{OrderId: anOrder.Id}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
	})

	t.Run("should return nil if already exists a transaction to the order", func(t *testing.T) {
		ctx := context.Background()

		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.TenantId)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		pixTransaction := transaction.NewPixTransaction(value_object.NewID(), anOrder.Id, "", faker.Word(), "", 10)
		err = applicationFactory.TransactionRepository.CreatePixTransaction(ctx, pixTransaction)
		require.NoError(t, err)

		useCase := New(applicationFactory)
		input := Input{OrderId: anOrder.Id}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
	})

	t.Run("should create a pix transaction", func(t *testing.T) {
		ctx := context.Background()

		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		anOrder, err := order.New(customerId, enums.Pix, enums.InApp, enums.Withdraw, faker.Word(), 0, 10, time.Now(), aMerchant.TenantId)
		require.NoError(t, err)

		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		require.NoError(t, err)

		pixOutput := &gateway.CreateQrCodePixOutput{
			Id:             value_object.NewID(),
			QrCode:         faker.Word(),
			ExpiresAt:      time.Now().Add(time.Hour),
			PaymentGateway: faker.Word(),
			QrCodeLink:     faker.URL(),
		}
		pixGateway.On("CreateQrCodePix", mock.Anything, mock.Anything).Return(pixOutput, nil).Once()

		useCase := New(applicationFactory)
		input := Input{OrderId: anOrder.Id}
		eventEmitter.On("Emit", ctx, mock.Anything).Return(nil)

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNumberOfCalls(t, "Emit", 1)
	})
}
