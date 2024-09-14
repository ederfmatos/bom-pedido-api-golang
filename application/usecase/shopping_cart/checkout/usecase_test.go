package checkout

import (
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CheckoutShoppingCart(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	shoppingCartRepository := applicationFactory.ShoppingCartRepository
	orderRepository := applicationFactory.OrderRepository
	productRepository := applicationFactory.ProductRepository
	useCase := New(applicationFactory)

	ctx := context.Background()
	t.Run("should return ShoppingCartEmptyError is shopping cart is empty", func(t *testing.T) {
		input := Input{CustomerId: value_object.NewID()}
		output, err := useCase.Execute(ctx, input)
		require.Nil(t, output)
		require.ErrorIs(t, err, errors.ShoppingCartEmptyError)

		shoppingCart := shopping_cart.New(input.CustomerId, faker.WORD)
		err = shoppingCartRepository.Upsert(ctx, shoppingCart)
		require.NoError(t, err)

		output, err = useCase.Execute(ctx, input)
		require.Nil(t, output)
		require.ErrorIs(t, err, errors.ShoppingCartEmptyError)
	})

	t.Run("should create a order", func(t *testing.T) {
		input := Input{
			CustomerId:      value_object.NewID(),
			PaymentMethod:   enums.CreditCard,
			DeliveryMode:    enums.Withdraw,
			PaymentMode:     enums.InReceiving,
			AddressId:       "",
			Payback:         0,
			CreditCardToken: "",
		}
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.NoError(t, err)

		err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
		require.NoError(t, err)

		aProduct, _ := product.New(faker.Name(), faker.Word(), 11.0, aMerchant.TenantId)
		err = productRepository.Create(ctx, aProduct)
		require.NoError(t, err)

		shoppingCart := shopping_cart.New(input.CustomerId, aMerchant.TenantId)
		err = shoppingCart.AddItem(aProduct, 1, "")
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
