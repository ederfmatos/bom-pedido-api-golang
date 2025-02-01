package category

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
)

func Test_CreateCategory(t *testing.T) {
	ctx := context.Background()
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("should return error if exists a category with the same name", func(t *testing.T) {
		category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())
		err := applicationFactory.ProductCategoryRepository.Create(ctx, category)
		require.NoError(t, err)

		useCase := NewCreateCategory(applicationFactory)
		input := CreateCategoryInput{
			TenantId:    category.TenantId,
			Name:        category.Name,
			Description: faker.Word(),
		}
		output, err := useCase.Execute(ctx, input)
		require.ErrorIs(t, err, CategoryWithSameNameError)
		require.Nil(t, output)

		input = CreateCategoryInput{
			TenantId:    value_object.NewTenantId(),
			Name:        category.Name,
			Description: faker.Word(),
		}
		output, err = useCase.Execute(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

	t.Run("should create a category", func(t *testing.T) {
		useCase := NewCreateCategory(applicationFactory)
		input := CreateCategoryInput{
			TenantId:    faker.Word(),
			Name:        faker.Name(),
			Description: faker.Word(),
		}
		output, err := useCase.Execute(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
	})
}
