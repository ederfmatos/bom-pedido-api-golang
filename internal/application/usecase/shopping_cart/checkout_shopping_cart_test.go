package shopping_cart

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CheckoutShoppingCart(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	shoppingCartRepository := applicationFactory.ShoppingCartRepository
	orderRepository := applicationFactory.OrderRepository
	productRepository := applicationFactory.ProductRepository
	useCase := NewCheckoutShoppingCart(applicationFactory)

	ctx := context.Background()
	t.Run("should return ShoppingCartEmptyError is shopping cart is empty", func(t *testing.T) {
		input := CheckoutShoppingCartInput{CustomerId: value_object.NewID()}
		output, err := useCase.Execute(ctx, input)
		require.Nil(t, output)
		require.ErrorIs(t, err, errors.ShoppingCartEmptyError)

		shoppingCart := entity.NewShoppingCart(input.CustomerId, faker.Word())
		err = shoppingCartRepository.Upsert(ctx, shoppingCart)
		require.NoError(t, err)

		output, err = useCase.Execute(ctx, input)
		require.Nil(t, output)
		require.ErrorIs(t, err, errors.ShoppingCartEmptyError)
	})

	t.Run("should create a order", func(t *testing.T) {
		input := CheckoutShoppingCartInput{
			CustomerId:      value_object.NewID(),
			PaymentMethod:   enums.CreditCard,
			DeliveryMode:    enums.Withdraw,
			PaymentMode:     enums.InReceiving,
			AddressId:       "",
			Payback:         0,
			CreditCardToken: "",
		}
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, merchant)
		require.NoError(t, err)

		product, _ := entity.NewProduct(faker.Name(), faker.Word(), 11.0, faker.Word(), merchant.TenantId)
		err = productRepository.Create(ctx, product)
		require.NoError(t, err)

		shoppingCart := entity.NewShoppingCart(input.CustomerId, merchant.TenantId)
		err = shoppingCart.AddItem(product, 1, "")
		require.NoError(t, err)

		err = shoppingCartRepository.Upsert(ctx, shoppingCart)
		require.NoError(t, err)

		output, err := useCase.Execute(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)

		order, err := orderRepository.FindById(ctx, output.Id)
		require.NoError(t, err)
		require.NotNil(t, order)

		require.Equal(t, shoppingCart.CustomerId, order.CustomerID)
		require.Equal(t, enums.PaymentMethodCreditCard, order.PaymentMethod)
		require.Equal(t, enums.DeliveryModeWithdraw, order.DeliveryMode)
		require.Equal(t, enums.PaymentModeInReceiving, order.PaymentMode)
		require.Equal(t, "", order.CreditCardToken)
		require.Equal(t, float64(0), order.Payback)
		require.Equal(t, int32(1), order.Code)
	})
}
