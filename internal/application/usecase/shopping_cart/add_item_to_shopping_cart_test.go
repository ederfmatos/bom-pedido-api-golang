package shopping_cart

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddItemToShoppingCartUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewAddItemToShoppingCart(applicationFactory)

	ctx := context.Background()
	t.Run("should return ProductNotFoundError", func(t *testing.T) {
		input := AddItemToShoppingCartInput{
			ProductId: value_object.NewID(),
			Quantity:  1,
		}
		err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.ProductNotFoundError)
	})

	t.Run("should return error is product is unavailable", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, faker.Word(), faker.Word())
		product.MarkUnAvailable()
		if err != nil {
			t.Fatalf("failed to restore product: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(ctx, product)

		input := AddItemToShoppingCartInput{
			CustomerId:  value_object.NewID(),
			ProductId:   product.Id,
			Quantity:    2,
			Observation: faker.Word(),
		}
		err = useCase.Execute(ctx, input)
		require.ErrorIs(t, err, errors.ProductUnAvailableError)

		shoppingCart, _ := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
		require.Nil(t, shoppingCart)
	})

	t.Run("should create a shopping cart with one item", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, faker.Word(), faker.Word())
		if err != nil {
			t.Fatalf("failed to restore product: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(ctx, product)

		input := AddItemToShoppingCartInput{
			CustomerId:  value_object.NewID(),
			ProductId:   product.Id,
			Quantity:    2,
			Observation: faker.Word(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		shoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
		require.NoError(t, err)
		require.NotNil(t, shoppingCart)
		require.Equal(t, 20.0, shoppingCart.GetPrice())
		require.Equal(t, 1, len(shoppingCart.Items))
		for _, item := range shoppingCart.Items {
			require.Equal(t, product.Id, item.ProductId)
			require.Equal(t, input.Quantity, item.Quantity)
			require.Equal(t, input.Observation, item.Observation)
			require.Equal(t, 10.0, item.Price)
			require.Equal(t, 20.0, item.GetTotalPrice())
		}
	})
}
