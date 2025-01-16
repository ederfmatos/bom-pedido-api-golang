package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"context"
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
				category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())

				product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				savedProduct, err := productRepository.FindById(ctx, product.Id)
				require.NoError(t, err)
				require.Nil(t, savedProduct)

				existsByName, err := productRepository.ExistsByNameAndTenantId(ctx, product.Name, faker.Word())
				require.NoError(t, err)
				require.False(t, existsByName)

				err = productRepository.Create(ctx, product)
				require.NoError(t, err)

				savedProduct, err = productRepository.FindById(ctx, product.Id)
				require.NoError(t, err)
				require.Equal(t, product.Id, savedProduct.Id)
				require.Equal(t, product.Name, savedProduct.Name)
				require.Equal(t, product.Description, savedProduct.Description)
				require.Equal(t, product.Price, savedProduct.Price)
				require.Equal(t, product.Status, savedProduct.Status)

				existsByName, err = productRepository.ExistsByNameAndTenantId(ctx, product.Name, product.TenantId)
				require.NoError(t, err)
				require.True(t, existsByName)
			})

			t.Run("should not create duplicated product", func(t *testing.T) {
				category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())

				product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				err = productRepository.Create(ctx, product)
				require.NoError(t, err)

				anotherProduct, err := entity.NewProduct(product.Name, product.Description, product.Price, product.CategoryId, product.TenantId)
				require.NoError(t, err)

				// TODO: Não permitir registro repetido
				//err = productRepository.Create(ctx, anotherProduct)
				//require.Error(t, err)

				anotherProduct.Name = faker.Name()
				err = productRepository.Create(ctx, anotherProduct)
				require.NoError(t, err)
			})

			t.Run("should update a product", func(t *testing.T) {
				category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())

				product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, category.Id, faker.Word())
				require.NoError(t, err)

				err = productRepository.Create(ctx, product)
				require.NoError(t, err)

				savedProduct, err := productRepository.FindById(ctx, product.Id)
				require.NoError(t, err)
				require.Equal(t, product, savedProduct)

				product.Name = faker.Name()
				product.Description = faker.Word()
				product.Price = 15.0
				err = productRepository.Update(ctx, product)
				require.NoError(t, err)

				savedProduct, err = productRepository.FindById(ctx, product.Id)
				require.NoError(t, err)
				require.Equal(t, product, savedProduct)
			})
		})
	}
}
