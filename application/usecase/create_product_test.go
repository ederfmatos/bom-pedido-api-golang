package usecase

import (
	"bom-pedido-api/domain/entity"
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
	useCase := NewCreateProductUseCase(applicationFactory)

	t.Run("should return ProductWithSameNameError error", func(t *testing.T) {
		input := CreateProductInput{
			Context:     context.Background(),
			Name:        faker.Name(),
			Description: faker.Word(),
			Price:       10.0,
		}
		product, err := entity.RestoreProduct(value_object.NewID(), input.Name, faker.Word(), 10.0, "ACTIVE")
		if err != nil {
			t.Fatalf("failed to restore product: %v", err)
		}
		_ = applicationFactory.ProductRepository.Create(input.Context, product)

		output, err := useCase.Execute(input)

		assert.ErrorIs(t, err, entity.ProductWithSameNameError)
		assert.Nil(t, output)
	})

	t.Run("should return an error is product is invalid", func(t *testing.T) {
		tests := []struct {
			name        string
			description string
			price       float64
			wantErr     error
		}{
			{name: "", description: "", price: 10, wantErr: errors.NewCompositeWithError(entity.ProductNameIsRequiredError)},
			{name: faker.Name(), description: "", price: 0, wantErr: errors.NewCompositeWithError(entity.ProductPriceIsRequiredError)},
			{name: faker.Name(), description: "", price: -1, wantErr: errors.NewCompositeWithError(entity.ProductPriceShouldPositiveError)},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("should return %s error", tt.wantErr.Error()), func(t *testing.T) {
				input := CreateProductInput{
					Context:     context.Background(),
					Name:        tt.name,
					Description: tt.description,
					Price:       tt.price,
				}

				output, err := useCase.Execute(input)

				assert.Equal(t, err.Error(), tt.wantErr.Error())
				assert.Nil(t, output)
			})
		}
	})

	t.Run("should create a product", func(t *testing.T) {
		input := CreateProductInput{
			Context:     context.Background(),
			Name:        faker.Name(),
			Description: faker.Word(),
			Price:       10.0,
		}

		output, err := useCase.Execute(input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
		savedProduct, _ := applicationFactory.ProductRepository.FindById(input.Context, output.ID)
		assert.NotNil(t, savedProduct)
		assert.Equal(t, input.Name, savedProduct.Name)
		assert.Equal(t, input.Description, savedProduct.Description)
		assert.Equal(t, input.Price, savedProduct.Price)
		assert.True(t, savedProduct.IsActive())
	})
}
