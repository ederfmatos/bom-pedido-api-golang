package order

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CloneOrder(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewCloneOrder(applicationFactory)

	t.Run("should return order not found", func(t *testing.T) {
		ctx := context.Background()
		input := CloneOrderInput{
			OrderId: value_object.NewID(),
		}
		err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.OrderNotFoundError)
	})

	t.Run("should clone an order", func(t *testing.T) {
		ctx := context.Background()
		customerId := value_object.NewID()
		order, err := entity.NewOrder(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 0, 0, time.Now(), faker.Word())
		require.NoError(t, err)
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, faker.Word(), faker.Word())
		require.NoError(t, err)
		err = order.AddProduct(product, 1, "observation")
		require.NoError(t, err)
		err = applicationFactory.OrderRepository.Create(ctx, order)
		require.NoError(t, err)

		input := CloneOrderInput{
			OrderId: order.Id,
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, order.CustomerID)
		require.NoError(t, err)
		require.NotNil(t, savedShoppingCart)
		require.Equal(t, len(savedShoppingCart.Items), 1)
		for _, shoppingCartItem := range savedShoppingCart.Items {
			require.Equal(t, shoppingCartItem.ProductId, product.Id)
			require.Equal(t, shoppingCartItem.Price, product.Price)
			require.Equal(t, shoppingCartItem.Quantity, 1)
			require.Equal(t, shoppingCartItem.Observation, "observation")
		}
	})
}
