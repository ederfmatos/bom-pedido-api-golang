package product

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	errors2 "errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		Product *Product
		Errors  []error
	}{
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: 0.0, Status: "RANDOM"},
			Errors:  []error{errors.ProductNameIsRequiredError, errors.ProductPriceIsRequiredError, errors.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: -1.0, Status: "RANDOM"},
			Errors:  []error{errors.ProductNameIsRequiredError, errors.ProductPriceShouldPositiveError, errors.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: 5.0, Status: "RANDOM"},
			Errors:  []error{errors.ProductNameIsRequiredError, errors.ProductInvalidProductStatusError},
		},
		{
			Product: &Product{Id: value_object.NewID(), Name: "", Price: -1.0, Status: "ACTIVE"},
			Errors:  []error{errors.ProductNameIsRequiredError, errors.ProductPriceShouldPositiveError},
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
			var composite *errors.CompositeError
			errors2.As(err, &composite)

			require.ErrorAs(t, err, &composite)
			require.Equal(t, len(test.Errors), len(composite.Errors))
			for index, err := range test.Errors {
				require.ErrorIs(t, composite.Errors[index], err)
			}
		})
	}
}
