package shopping_cart

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewDeleteShoppingCart(applicationFactory)
	customerId := value_object.NewID()
	shoppingCart := entity.NewShoppingCart(customerId, faker.Word())
	input := DeleteShoppingCartInput{
		CustomerId: customerId,
	}
	ctx := context.Background()
	err := applicationFactory.ShoppingCartRepository.Upsert(ctx, shoppingCart)
	require.NoError(t, err)
	err = useCase.Execute(ctx, input)
	require.NoError(t, err)
	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, customerId)
	require.NoError(t, err)
	require.Nil(t, savedShoppingCart)
}
