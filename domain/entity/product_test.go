package entity

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	errors2 "errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		Product *Product
		Errors  []error
	}{
		{
			Product: &Product{ID: value_object.NewID(), Name: "", Price: 0.0, Status: "RANDOM"},
			Errors:  []error{ProductNameIsRequiredError, ProductPriceIsRequiredError, InvalidProductStatusError},
		},
		{
			Product: &Product{ID: value_object.NewID(), Name: "", Price: -1.0, Status: "RANDOM"},
			Errors:  []error{ProductNameIsRequiredError, ProductPriceShouldPositiveError, InvalidProductStatusError},
		},
		{
			Product: &Product{ID: value_object.NewID(), Name: "", Price: 5.0, Status: "RANDOM"},
			Errors:  []error{ProductNameIsRequiredError, InvalidProductStatusError},
		},
		{
			Product: &Product{ID: value_object.NewID(), Name: "", Price: -1.0, Status: "ACTIVE"},
			Errors:  []error{ProductNameIsRequiredError, ProductPriceShouldPositiveError},
		},
		{
			Product: &Product{ID: value_object.NewID(), Name: faker.Name(), Price: 10.0, Status: "ACTIVE"},
			Errors:  []error{},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("should return errors %v", test.Errors), func(t *testing.T) {
			err := test.Product.Validate()
			if len(test.Errors) == 0 {
				assert.Nil(t, err)
				return
			}

			assert.Error(t, err)
			var composite *errors.CompositeError
			errors2.As(err, &composite)

			assert.ErrorAs(t, err, &composite)
			assert.Equal(t, len(test.Errors), len(composite.Errors))
			for index, err := range test.Errors {
				assert.ErrorIs(t, composite.Errors[index], err)
			}
		})
	}
}
