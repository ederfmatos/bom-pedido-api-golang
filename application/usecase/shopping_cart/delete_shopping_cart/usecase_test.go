package delete_shopping_cart

import (
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)
	customerId := value_object.NewID()
	aShoppingCart := shopping_cart.New(customerId, faker.WORD)
	input := Input{
		CustomerId: customerId,
	}
	ctx := context.Background()
	err := applicationFactory.ShoppingCartRepository.Upsert(ctx, aShoppingCart)
	assert.NoError(t, err)
	err = useCase.Execute(ctx, input)
	assert.NoError(t, err)
	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, customerId)
	assert.NoError(t, err)
	assert.Nil(t, savedShoppingCart)
}
