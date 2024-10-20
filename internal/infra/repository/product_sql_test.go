package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ProductSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	productRepository := NewDefaultProductRepository(sqlConnection)
	categoryRepository := NewDefaultProductCategoryRepository(sqlConnection)
	runProductTests(t, productRepository, categoryRepository)
}

func Test_ProductMemoryRepository(t *testing.T) {
	productRepository := NewProductMemoryRepository()
	categoryRepository := NewProductCategoryMemoryRepository()
	runProductTests(t, productRepository, categoryRepository)
}

func runProductTests(t *testing.T, repository repository.ProductRepository, categoryRepository repository.ProductCategoryRepository) {
	ctx := context.Background()

	t.Run("should create a product", func(t *testing.T) {
		category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
		err := categoryRepository.Create(ctx, category)
		require.NoError(t, err)

		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
		require.NoError(t, err)

		savedProduct, err := repository.FindById(ctx, aProduct.Id)
		require.NoError(t, err)
		require.Nil(t, savedProduct)

		existsByName, err := repository.ExistsByNameAndTenantId(ctx, aProduct.Name, faker.WORD)
		require.NoError(t, err)
		require.False(t, existsByName)

		err = repository.Create(ctx, aProduct)
		require.NoError(t, err)

		savedProduct, err = repository.FindById(ctx, aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct.Id, savedProduct.Id)
		require.Equal(t, aProduct.Name, savedProduct.Name)
		require.Equal(t, aProduct.Description, savedProduct.Description)
		require.Equal(t, aProduct.Price, savedProduct.Price)
		require.Equal(t, aProduct.Status, savedProduct.Status)

		existsByName, err = repository.ExistsByNameAndTenantId(ctx, aProduct.Name, aProduct.TenantId)
		require.NoError(t, err)
		require.True(t, existsByName)
	})

	t.Run("should not create duplicated product", func(t *testing.T) {
		category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
		err := categoryRepository.Create(ctx, category)
		require.NoError(t, err)

		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
		require.NoError(t, err)

		err = repository.Create(ctx, aProduct)
		require.NoError(t, err)

		anotherProduct, err := product.New(aProduct.Name, aProduct.Description, aProduct.Price, aProduct.CategoryId, aProduct.TenantId)
		require.NoError(t, err)

		err = repository.Create(ctx, anotherProduct)
		require.Error(t, err)

		anotherProduct.Name = faker.Name()
		err = repository.Create(ctx, anotherProduct)
		require.NoError(t, err)
	})

	t.Run("should update a product", func(t *testing.T) {
		category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())
		err := categoryRepository.Create(ctx, category)
		require.NoError(t, err)

		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
		require.NoError(t, err)

		err = repository.Create(ctx, aProduct)
		require.NoError(t, err)

		savedProduct, err := repository.FindById(ctx, aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct, savedProduct)

		aProduct.Name = faker.Name()
		aProduct.Description = faker.Word()
		aProduct.Price = 15.0
		err = repository.Update(ctx, aProduct)
		require.NoError(t, err)

		savedProduct, err = repository.FindById(ctx, aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct, savedProduct)
	})
}
