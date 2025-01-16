package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
)

func TestCategoriesRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.ProductCategoryRepository{
		"CategoriesMemoryRepository": NewCategoriesMemoryRepository(),
		"CategoriesMongoRepository":  NewCategoriesMongoRepository(container.MongoDatabase()),
	}

	for name, categoryRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			category := entity.NewCategory(faker.Name(), faker.Word(), value_object.NewTenantId())

			existsById, err := categoryRepository.ExistsById(ctx, category.Id)
			require.NoError(t, err)
			require.False(t, existsById)

			existsByNameAndTenantId, err := categoryRepository.ExistsByNameAndTenantId(ctx, category.Name, category.TenantId)
			require.NoError(t, err)
			require.False(t, existsByNameAndTenantId)

			err = categoryRepository.Create(ctx, category)
			require.NoError(t, err)

			existsById, err = categoryRepository.ExistsById(ctx, category.Id)
			require.NoError(t, err)
			require.True(t, existsById)

			existsByNameAndTenantId, err = categoryRepository.ExistsByNameAndTenantId(ctx, category.Name, category.TenantId)
			require.NoError(t, err)
			require.True(t, existsByNameAndTenantId)
		})
	}
}
