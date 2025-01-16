package entity

import (
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
)

type ProductStatus string

const (
	ProductActive      = ProductStatus("ACTIVE")
	ProductInactive    = ProductStatus("INACTIVE")
	ProductUnAvailable = ProductStatus("UNAVAILABLE")
)

type Product struct {
	Id          string        `bson:"id"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Price       float64       `bson:"price"`
	Status      ProductStatus `bson:"status"`
	TenantId    string        `bson:"tenantId"`
	CategoryId  string        `bson:"categoryId"`
}

func NewProduct(name, description string, price float64, categoryId, tenantId string) (*Product, error) {
	product := &Product{
		Id:          value_object.NewID(),
		Name:        name,
		Price:       price,
		Description: description,
		Status:      ProductActive,
		TenantId:    tenantId,
		CategoryId:  categoryId,
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
	if product.Status != ProductActive && product.Status != ProductInactive {
		compositeError.Append(errors.ProductInvalidProductStatusError)
	}
	return compositeError.AsError()
}

func (product *Product) IsActive() bool {
	return product.Status == ProductActive
}

func (product *Product) IsUnAvailable() bool {
	return product.Status == ProductUnAvailable
}

func (product *Product) MarkUnAvailable() {
	product.Status = ProductUnAvailable
}
