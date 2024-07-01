package entity

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
)

type ProductStatus string

var (
	ProductStatusActive      = ProductStatus("ACTIVE")
	ProductStatusInactive    = ProductStatus("INACTIVE")
	ProductStatusUnAvailable = ProductStatus("UNAVAILABLE")
)

var (
	ProductNameIsRequiredError      = errors.New("product name is required")
	ProductPriceIsRequiredError     = errors.New("product price is required")
	ProductPriceShouldPositiveError = errors.New("product price should positive")
	InvalidProductStatusError       = errors.New("invalid product status")
	ProductWithSameNameError        = errors.New("product with this name already exists")
	ProductUnAvailableError         = errors.New("Produto indisponível")
	ProductNotFoundError            = errors.New("Produto não encontrado")
)

type Product struct {
	Id          string
	Name        string
	Description string
	Price       float64
	Status      ProductStatus
}

func NewProduct(name, description string, price float64) (*Product, error) {
	product := &Product{
		Id:          value_object.NewID(),
		Name:        name,
		Price:       price,
		Description: description,
		Status:      ProductStatusActive,
	}
	return product, product.Validate()
}

func RestoreProduct(id, name, description string, price float64, status string) (*Product, error) {
	product := &Product{
		Id:          id,
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
		compositeError.Append(ProductNameIsRequiredError)
	}
	if product.Price == 0.0 {
		compositeError.Append(ProductPriceIsRequiredError)
	}
	if product.Price < 0.0 {
		compositeError.Append(ProductPriceShouldPositiveError)
	}
	if product.Status != ProductStatusActive && product.Status != ProductStatusInactive {
		compositeError.Append(InvalidProductStatusError)
	}
	return compositeError.AsError()
}

func (product *Product) IsActive() bool {
	return product.Status == ProductStatusActive
}

func (product *Product) IsUnAvailable() bool {
	return product.Status == ProductStatusUnAvailable
}

func (product *Product) MarkUnAvailable() {
	product.Status = ProductStatusUnAvailable
}
