package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
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

			aCategory := product.NewCategory(faker.Name(), faker.Word(), value_object.NewTenantId())

			existsById, err := categoryRepository.ExistsById(ctx, aCategory.Id)
			require.NoError(t, err)
			require.False(t, existsById)

			existsByNameAndTenantId, err := categoryRepository.ExistsByNameAndTenantId(ctx, aCategory.Name, aCategory.TenantId)
			require.NoError(t, err)
			require.False(t, existsByNameAndTenantId)

			err = categoryRepository.Create(ctx, aCategory)
			require.NoError(t, err)

			existsById, err = categoryRepository.ExistsById(ctx, aCategory.Id)
			require.NoError(t, err)
			require.True(t, existsById)

			existsByNameAndTenantId, err = categoryRepository.ExistsByNameAndTenantId(ctx, aCategory.Name, aCategory.TenantId)
			require.NoError(t, err)
			require.True(t, existsByNameAndTenantId)
		})
	}
}
