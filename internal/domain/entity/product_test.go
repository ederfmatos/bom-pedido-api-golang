package entity

import (
	domainError "bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"errors"
	"fmt"
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		Product *Product
		Errors  []error
	}{
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: 0.0, Status: "RANDOM"},
			Errors:  []error{domainError.ProductNameIsRequiredError, domainError.ProductPriceIsRequiredError, domainError.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: -1.0, Status: "RANDOM"},
			Errors:  []error{domainError.ProductNameIsRequiredError, domainError.ProductPriceShouldPositiveError, domainError.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: 5.0, Status: "RANDOM"},
			Errors:  []error{domainError.ProductNameIsRequiredError, domainError.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: -1.0, Status: "ACTIVE"},
			Errors:  []error{domainError.ProductNameIsRequiredError, domainError.ProductPriceShouldPositiveError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: faker.Name(), Price: 10.0, Status: "ACTIVE"},
			Errors:  []error{},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("should return errors %v", test.Errors), func(t *testing.T) {
			err := test.Product.Validate()
			if len(test.Errors) == 0 {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			var composite *domainError.CompositeError
			errors.As(err, &composite)

			require.ErrorAs(t, err, &composite)
			require.Equal(t, len(test.Errors), len(composite.Errors))
			for index, err := range test.Errors {
				require.ErrorIs(t, composite.Errors[index], err)
			}
		})
	}
}
