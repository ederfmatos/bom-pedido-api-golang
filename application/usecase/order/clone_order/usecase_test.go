package clone_order

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_UseCase(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := Input{
			OrderId: value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		assert.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should clone an order", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		anOrder, err := order.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, time.Now(), faker.Word())
		assert.NoError(t, err)
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, faker.Word())
		assert.NoError(t, err)
		err = anOrder.AddProduct(aProduct, 1, "observation")
		assert.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, anOrder)
		assert.NoError(t, err)

		input := Input{
			OrderId: anOrder.Id,
		}
		err = useCase.Execute(ctx, input)
		assert.NoError(t, err)

		savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, anOrder.CustomerID)
		assert.NoError(t, err)
		assert.NotNil(t, savedShoppingCart)
		assert.Equal(t, len(savedShoppingCart.Items), 1)
		shoppingCartItem := savedShoppingCart.Items[0]
		assert.Equal(t, shoppingCartItem.ProductId, aProduct.Id)
		assert.Equal(t, shoppingCartItem.Price, aProduct.Price)
		assert.Equal(t, shoppingCartItem.Quantity, 1)
		assert.Equal(t, shoppingCartItem.Observation, "observation")
	})
}
