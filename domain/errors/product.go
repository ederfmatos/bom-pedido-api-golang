package errors

import "errors"

var (
	ProductNameIsRequired      = errors.New("product name is required")
	ProductPriceIsRequired     = errors.New("product price is required")
	ProductPriceShouldPositive = errors.New("product price should positive")
	InvalidProductStatus       = errors.New("invalid product status")
	ProductWithSameName        = errors.New("product with this name already exists")
	ProductUnAvailable         = errors.New("Produto indisponível")
	ProductNotFoundError       = errors.New("Produto não encontrado")
)
