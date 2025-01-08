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

func Test_ProductRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.ProductRepository{
		"ProductMemoryRepository": NewProductMemoryRepository(),
		"ProductMongoRepository":  NewProductMongoRepository(container.MongoDatabase()),
	}

	for name, productRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			t.Run("should create a product", func(t *testing.T) {
				category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())

				aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				savedProduct, err := productRepository.FindById(ctx, aProduct.Id)
				require.NoError(t, err)
				require.Nil(t, savedProduct)

				existsByName, err := productRepository.ExistsByNameAndTenantId(ctx, aProduct.Name, faker.WORD)
				require.NoError(t, err)
				require.False(t, existsByName)

				err = productRepository.Create(ctx, aProduct)
				require.NoError(t, err)

				savedProduct, err = productRepository.FindById(ctx, aProduct.Id)
				require.NoError(t, err)
				require.Equal(t, aProduct.Id, savedProduct.Id)
				require.Equal(t, aProduct.Name, savedProduct.Name)
				require.Equal(t, aProduct.Description, savedProduct.Description)
				require.Equal(t, aProduct.Price, savedProduct.Price)
				require.Equal(t, aProduct.Status, savedProduct.Status)

				existsByName, err = productRepository.ExistsByNameAndTenantId(ctx, aProduct.Name, aProduct.TenantId)
				require.NoError(t, err)
				require.True(t, existsByName)
			})

			t.Run("should not create duplicated product", func(t *testing.T) {
				category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())

				aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				err = productRepository.Create(ctx, aProduct)
				require.NoError(t, err)

				anotherProduct, err := product.New(aProduct.Name, aProduct.Description, aProduct.Price, aProduct.CategoryId, aProduct.TenantId)
				require.NoError(t, err)

				// TODO: Não permitir registro repetido
				//err = productRepository.Create(ctx, anotherProduct)
				//require.Error(t, err)

				anotherProduct.Name = faker.Name()
				err = productRepository.Create(ctx, anotherProduct)
				require.NoError(t, err)
			})

			t.Run("should update a product", func(t *testing.T) {
				category := product.NewCategory(faker.Name(), faker.Word(), faker.Word())

				aProduct, err := product.New(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				err = productRepository.Create(ctx, aProduct)
				require.NoError(t, err)

				savedProduct, err := productRepository.FindById(ctx, aProduct.Id)
				require.NoError(t, err)
				require.Equal(t, aProduct, savedProduct)

				aProduct.Name = faker.Name()
				aProduct.Description = faker.Word()
				aProduct.Price = 15.0
				err = productRepository.Update(ctx, aProduct)
				require.NoError(t, err)

				savedProduct, err = productRepository.FindById(ctx, aProduct.Id)
				require.NoError(t, err)
				require.Equal(t, aProduct, savedProduct)
			})
		})
	}
}
