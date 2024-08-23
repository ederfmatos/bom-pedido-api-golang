package create_product

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProductUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := New(applicationFactory)
	ctx := context.TODO()

	t.Run("should return ProductWithSameNameError error", func(t *testing.T) {
		input := Input{
			Name:        faker.Name(),
			Description: faker.Word(),
			Price:       10.0,
			TenantId:    faker.Word(),
		}
		aProduct, err := product.Restore(value_object.NewID(), input.Name, faker.Word(), 10.0, "ACTIVE", input.TenantId)
		if err != nil {
			t.Fatalf("failed to restore aProduct: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(ctx, aProduct)

		output, err := useCase.Execute(ctx, input)

		assert.ErrorIs(t, err, errors.ProductWithSameNameError)
		assert.Nil(t, output)
	})

	t.Run("should return an error is product is invalid", func(t *testing.T) {
		tests := []struct {
			name        string
			description string
			price       float64
			wantErr     error
		}{
			{name: "", description: "", price: 10, wantErr: errors.NewCompositeWithError(errors.ProductNameIsRequiredError)},
			{name: faker.Name(), description: "", price: 0, wantErr: errors.NewCompositeWithError(errors.ProductPriceIsRequiredError)},
			{name: faker.Name(), description: "", price: -1, wantErr: errors.NewCompositeWithError(errors.ProductPriceShouldPositiveError)},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("should return %s error", tt.wantErr.Error()), func(t *testing.T) {
				input := Input{
					Name:        tt.name,
					Description: tt.description,
					Price:       tt.price,
					TenantId:    faker.Word(),
				}

				output, err := useCase.Execute(ctx, input)

				assert.Equal(t, err.Error(), tt.wantErr.Error())
				assert.Nil(t, output)
			})
		}
	})

	t.Run("should create a product", func(t *testing.T) {
		input := Input{
			Name:        faker.Name(),
			Description: faker.Word(),
			Price:       10.0,
			TenantId:    faker.Word(),
		}

		output, err := useCase.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.Id)
		savedProduct, _ := applicationFactory.ProductRepository.FindById(ctx, output.Id)
		assert.NotNil(t, savedProduct)
		assert.Equal(t, input.Name, savedProduct.Name)
		assert.Equal(t, input.Description, savedProduct.Description)
		assert.Equal(t, input.Price, savedProduct.Price)
		assert.Equal(t, input.TenantId, savedProduct.TenantId)
		assert.True(t, savedProduct.IsActive())
	})
}
