package usecase

import (
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddItemToShoppingCartUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewAddItemToShoppingCartUseCase(applicationFactory)

	t.Run("should return ProductNotFoundError", func(t *testing.T) {
		input := AddItemToShoppingCartInput{
			Context:   context.Background(),
			ProductId: value_object.NewID(),
		}
		err := useCase.Execute(input)
		assert.ErrorIs(t, err, entity.ProductNotFoundError)
	})

	t.Run("should return error is product is unavailable", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0)
		product.MarkUnAvailable()
		if err != nil {
			t.Fatalf("failed to restore product: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(context.Background(), product)

		input := AddItemToShoppingCartInput{
			Context:     context.Background(),
			CustomerId:  value_object.NewID(),
			ProductId:   product.Id,
			Quantity:    2,
			Observation: faker.Word(),
		}
		err = useCase.Execute(input)
		assert.ErrorIs(t, err, entity.ProductUnAvailableError)

		shoppingCart, _ := applicationFactory.ShoppingCartRepository.FindByCustomerId(input.Context, input.CustomerId)
		assert.Nil(t, shoppingCart)
	})

	t.Run("should create a shopping cart with one item", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0)
		if err != nil {
			t.Fatalf("failed to restore product: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(context.Background(), product)

		input := AddItemToShoppingCartInput{
			Context:     context.Background(),
			CustomerId:  value_object.NewID(),
			ProductId:   product.Id,
			Quantity:    2,
			Observation: faker.Word(),
		}
		err = useCase.Execute(input)
		assert.NoError(t, err)

		shoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(input.Context, input.CustomerId)
		assert.NoError(t, err)
		assert.NotNil(t, shoppingCart)
		assert.Equal(t, 20.0, shoppingCart.GetPrice())
		items := shoppingCart.GetItems()
		assert.Equal(t, 1, len(items))
		item := items[0]
		assert.Equal(t, product.Id, item.ProductId)
		assert.Equal(t, input.Quantity, item.Quantity)
		assert.Equal(t, input.Observation, item.Observation)
		assert.Equal(t, 10.0, item.Price)
		assert.Equal(t, 20.0, item.GetTotalPrice())
	})
}
