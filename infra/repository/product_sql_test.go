package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ProductSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	productRepository := NewDefaultProductRepository(sqlConnection)
	runProductTests(t, productRepository)
}

func Test_ProductMemoryRepository(t *testing.T) {
	productRepository := NewProductMemoryRepository()
	runProductTests(t, productRepository)
}

func runProductTests(t *testing.T, repository repository.ProductRepository) {
	t.Run("should create a product", func(t *testing.T) {
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, faker.WORD)
		require.NoError(t, err)

		savedProduct, err := repository.FindById(context.Background(), aProduct.Id)
		require.NoError(t, err)
		require.Nil(t, savedProduct)

		existsByName, err := repository.ExistsByNameAndTenantId(context.Background(), aProduct.Name, faker.WORD)
		require.NoError(t, err)
		require.False(t, existsByName)

		err = repository.Create(context.Background(), aProduct)
		require.NoError(t, err)

		savedProduct, err = repository.FindById(context.Background(), aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct.Id, savedProduct.Id)
		require.Equal(t, aProduct.Name, savedProduct.Name)
		require.Equal(t, aProduct.Description, savedProduct.Description)
		require.Equal(t, aProduct.Price, savedProduct.Price)
		require.Equal(t, aProduct.Status, savedProduct.Status)

		existsByName, err = repository.ExistsByNameAndTenantId(context.Background(), aProduct.Name, faker.WORD)
		require.NoError(t, err)
		require.True(t, existsByName)
	})

	t.Run("should not create duplicated product", func(t *testing.T) {
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, faker.WORD)
		require.NoError(t, err)

		ctx := context.Background()
		err = repository.Create(ctx, aProduct)
		require.NoError(t, err)

		anotherProduct, err := product.New(aProduct.Name, aProduct.Description, aProduct.Price, aProduct.TenantId)
		require.NoError(t, err)

		err = repository.Create(ctx, anotherProduct)
		require.Error(t, err)

		anotherProduct.Name = faker.Name()
		err = repository.Create(ctx, anotherProduct)
		require.NoError(t, err)
	})

	t.Run("should update a product", func(t *testing.T) {
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, faker.WORD)
		require.NoError(t, err)

		err = repository.Create(context.Background(), aProduct)
		require.NoError(t, err)

		savedProduct, err := repository.FindById(context.Background(), aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct, savedProduct)

		aProduct.Name = faker.Name()
		aProduct.Description = faker.Word()
		aProduct.Price = 15.0
		err = repository.Update(context.Background(), aProduct)
		require.NoError(t, err)

		savedProduct, err = repository.FindById(context.Background(), aProduct.Id)
		require.NoError(t, err)
		require.Equal(t, aProduct, savedProduct)
	})
}
