package errors

var (
	ShoppingCartEmptyError       = New("Your shopping cart is empty")
	CardTokenIsRequiredError     = New("The card token is required")
	PaybackShouldBePositiveError = New("The payback should be positive")
)
