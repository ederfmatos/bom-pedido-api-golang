package product

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
)

type Status string

const (
	Active      = Status("ACTIVE")
	Inactive    = Status("INACTIVE")
	UnAvailable = Status("UNAVAILABLE")
)

type Product struct {
	Id          string
	Name        string
	Description string
	Price       float64
	Status      Status
}

func New(name, description string, price float64) (*Product, error) {
	product := &Product{
		Id:          value_object.NewID(),
		Name:        name,
		Price:       price,
		Description: description,
		Status:      Active,
	}
	return product, product.Validate()
}

func Restore(id, name, description string, price float64, status string) (*Product, error) {
	product := &Product{
		Id:          id,
		Name:        name,
		Price:       price,
		Description: description,
		Status:      Status(status),
	}
	return product, product.Validate()
}

func (product *Product) Validate() error {
	compositeError := errors.NewCompositeError()
	if product.Name == "" {
		compositeError.Append(errors.ProductNameIsRequiredError)
	}
	if product.Price == 0.0 {
		compositeError.Append(errors.ProductPriceIsRequiredError)
	}
	if product.Price < 0.0 {
		compositeError.Append(errors.ProductPriceShouldPositiveError)
	}
	if product.Status != Active && product.Status != Inactive {
		compositeError.Append(errors.ProductInvalidProductStatusError)
	}
	return compositeError.AsError()
}

func (product *Product) IsActive() bool {
	return product.Status == Active
}

func (product *Product) IsUnAvailable() bool {
	return product.Status == UnAvailable
}

func (product *Product) MarkUnAvailable() {
	product.Status = UnAvailable
}
