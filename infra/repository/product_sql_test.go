package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
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
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0)
		assert.NoError(t, err)

		savedProduct, err := repository.FindById(context.Background(), aProduct.Id)
		assert.NoError(t, err)
		assert.Nil(t, savedProduct)

		existsByName, err := repository.ExistsByName(context.Background(), aProduct.Name)
		assert.NoError(t, err)
		assert.False(t, existsByName)

		err = repository.Create(context.Background(), aProduct)
		assert.NoError(t, err)

		savedProduct, err = repository.FindById(context.Background(), aProduct.Id)
		assert.NoError(t, err)
		assert.Equal(t, aProduct.Id, savedProduct.Id)
		assert.Equal(t, aProduct.Name, savedProduct.Name)
		assert.Equal(t, aProduct.Description, savedProduct.Description)
		assert.Equal(t, aProduct.Price, savedProduct.Price)
		assert.Equal(t, aProduct.Status, savedProduct.Status)

		existsByName, err = repository.ExistsByName(context.Background(), aProduct.Name)
		assert.NoError(t, err)
		assert.True(t, existsByName)
	})

	t.Run("should not create duplicated product", func(t *testing.T) {
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0)
		assert.NoError(t, err)

		err = repository.Create(context.Background(), aProduct)
		assert.NoError(t, err)

		anotherProduct, err := product.New(aProduct.Name, aProduct.Description, aProduct.Price)
		assert.NoError(t, err)

		err = repository.Create(context.Background(), anotherProduct)
		assert.Error(t, err)

		anotherProduct.Name = faker.Name()
		err = repository.Create(context.Background(), anotherProduct)
		assert.NoError(t, err)
	})

	t.Run("should update a product", func(t *testing.T) {
		aProduct, err := product.New(faker.Name(), faker.Word(), 10.0)
		assert.NoError(t, err)

		err = repository.Create(context.Background(), aProduct)
		assert.NoError(t, err)

		savedProduct, err := repository.FindById(context.Background(), aProduct.Id)
		assert.NoError(t, err)
		assert.Equal(t, aProduct, savedProduct)

		aProduct.Name = faker.Name()
		aProduct.Description = faker.Word()
		aProduct.Price = 15.0
		err = repository.Update(context.Background(), aProduct)
		assert.NoError(t, err)

		savedProduct, err = repository.FindById(context.Background(), aProduct.Id)
		assert.NoError(t, err)
		assert.Equal(t, aProduct, savedProduct)
	})
}
