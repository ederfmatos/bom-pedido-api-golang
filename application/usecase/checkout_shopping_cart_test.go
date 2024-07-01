package usecase

import (
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/enums"
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
	useCase := NewCheckoutShoppingCartUseCase(applicationFactory)

	t.Run("should return ShoppingCartEmptyError is shopping cart is empty", func(t *testing.T) {
		input := CheckoutShoppingCartInput{Context: context.Background(), CustomerId: value_object.NewID()}
		output, err := useCase.Execute(input)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, entity.ShoppingCartEmptyError)

		shoppingCart := entity.NewShoppingCart(input.CustomerId)
		err = shoppingCartRepository.Upsert(input.Context, shoppingCart)
		assert.NoError(t, err)

		output, err = useCase.Execute(input)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, entity.ShoppingCartEmptyError)
	})

	t.Run("should create a order", func(t *testing.T) {
		input := CheckoutShoppingCartInput{
			Context:         context.Background(),
			CustomerId:      value_object.NewID(),
			PaymentMethod:   enums.CreditCard,
			DeliveryMode:    enums.Withdraw,
			PaymentMode:     enums.InReceiving,
			AddressId:       "",
			Change:          0,
			CreditCardToken: "",
		}
		product, _ := entity.NewProduct(faker.Name(), faker.Word(), 11.0)
		err := productRepository.Create(input.Context, product)
		assert.NoError(t, err)

		shoppingCart := entity.NewShoppingCart(input.CustomerId)
		err = shoppingCart.AddItem(product, 1, "")
		assert.NoError(t, err)

		err = shoppingCartRepository.Upsert(input.Context, shoppingCart)
		assert.NoError(t, err)

		output, err := useCase.Execute(input)
		assert.NoError(t, err)
		assert.NotNil(t, output)

		order, err := orderRepository.FindById(input.Context, output.Id)
		assert.NoError(t, err)
		assert.NotNil(t, order)

		assert.Equal(t, shoppingCart.CustomerId, order.CustomerID)
		assert.Equal(t, enums.PaymentMethodCreditCard, order.PaymentMethod)
		assert.Equal(t, enums.DeliveryModeWithdraw, order.DeliveryMode)
		assert.Equal(t, enums.PaymentModeInReceiving, order.PaymentMode)
		assert.Equal(t, "", order.CreditCardToken)
		assert.Equal(t, float64(0), order.Change)
		assert.Equal(t, int32(1), order.Code)
	})
}
