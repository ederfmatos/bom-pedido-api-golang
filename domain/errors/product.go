package errors

var (
	ProductNotFoundError    = New("Produto não encontrado")
	ProductUnAvailableError = New("Produto indisponível")

	ProductNameIsRequiredError       = New("product name is required")
	ProductPriceIsRequiredError      = New("product price is required")
	ProductPriceShouldPositiveError  = New("product price should positive")
	ProductInvalidProductStatusError = New("invalid product status")
	ProductWithSameNameError         = New("product with this name already exists")
)
