package create_category

import (
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CreateCategory(t *testing.T) {
	ctx := context.Background()
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("should return error if exists a category with the same name", func(t *testing.T) {
		category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
		err := applicationFactory.ProductCategoryRepository.Create(ctx, category)
		require.NoError(t, err)

		useCase := New(applicationFactory)
		input := Input{
			TenantId:    category.TenantId,
			Name:        category.Name,
			Description: faker.Word(),
		}
		err = useCase.Execute(ctx, input)
		require.ErrorIs(t, err, CategoryWithSameNameError)

		input = Input{
			TenantId:    value_object.NewTenantId(),
			Name:        category.Name,
			Description: faker.Word(),
		}
		err = useCase.Execute(ctx, input)
		require.NoError(t, err)
	})

	t.Run("should create a category", func(t *testing.T) {
		useCase := New(applicationFactory)
		input := Input{
			TenantId:    faker.Word(),
			Name:        faker.Name(),
			Description: faker.Word(),
		}
		err := useCase.Execute(ctx, input)
		require.NoError(t, err)
	})
}
