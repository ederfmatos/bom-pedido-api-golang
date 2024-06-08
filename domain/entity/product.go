package entity

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
)

type ProductStatus string

var (
	ProductStatusActive   = ProductStatus("ACTIVE")
	ProductStatusInactive = ProductStatus("INACTIVE")
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Status      ProductStatus
}

func NewProduct(name, description string, price float64) (*Product, error) {
	product := &Product{
		ID:          value_object.NewID(),
		Name:        name,
		Price:       price,
		Description: description,
		Status:      ProductStatusActive,
	}
	return product, product.Validate()
}

func RestoreProduct(ID, name, description string, price float64, status string) (*Product, error) {
	product := &Product{
		ID:          ID,
		Name:        name,
		Price:       price,
		Description: description,
		Status:      ProductStatus(status),
	}
	return product, product.Validate()
}

func (product *Product) Validate() error {
	compositeError := errors.NewCompositeError()
	if product.Name == "" {
		compositeError.Append(errors.ProductNameIsRequired)
	}
	if product.Price == 0.0 {
		compositeError.Append(errors.ProductPriceIsRequired)
	}
	if product.Price < 0.0 {
		compositeError.Append(errors.ProductPriceShouldPositive)
	}
	if product.Status != ProductStatusActive && product.Status != ProductStatusInactive {
		compositeError.Append(errors.InvalidProductStatus)
	}
	return compositeError.AsError()
}

func (product *Product) IsActive() bool {
	return product.Status == ProductStatusActive
}
