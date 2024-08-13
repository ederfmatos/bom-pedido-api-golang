package checkout

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CheckoutShoppingCart(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	shoppingCartRepository := applicationFactory.ShoppingCartRepository
	orderRepository := applicationFactory.OrderRepository
	productRepository := applicationFactory.ProductRepository
	useCase := New(applicationFactory)

	ctx := context.TODO()
	t.Run("should return ShoppingCartEmptyError is shopping cart is empty", func(t *testing.T) {
		input := Input{CustomerId: value_object.NewID()}
		output, err := useCase.Execute(ctx, input)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, errors.ShoppingCartEmptyError)

		shoppingCart := shopping_cart.New(input.CustomerId)
		err = shoppingCartRepository.Upsert(ctx, shoppingCart)
		assert.NoError(t, err)

		output, err = useCase.Execute(ctx, input)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, errors.ShoppingCartEmptyError)
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
		product, _ := product.New(faker.Name(), faker.Word(), 11.0)
		err := productRepository.Create(ctx, product)
		assert.NoError(t, err)

		shoppingCart := shopping_cart.New(input.CustomerId)
		err = shoppingCart.AddItem(product, 1, "")
		assert.NoError(t, err)

		err = shoppingCartRepository.Upsert(ctx, shoppingCart)
		assert.NoError(t, err)

		output, err := useCase.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, output)

		order, err := orderRepository.FindById(ctx, output.Id)
		assert.NoError(t, err)
		assert.NotNil(t, order)

		assert.Equal(t, shoppingCart.CustomerId, order.CustomerID)
		assert.Equal(t, enums.PaymentMethodCreditCard, order.PaymentMethod)
		assert.Equal(t, enums.DeliveryModeWithdraw, order.DeliveryMode)
		assert.Equal(t, enums.PaymentModeInReceiving, order.PaymentMode)
		assert.Equal(t, "", order.CreditCardToken)
		assert.Equal(t, float64(0), order.Payback)
		assert.Equal(t, int32(1), order.Code)
	})
}
