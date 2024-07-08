package errors

var (
	ShoppingCartEmptyError      = New("Your shopping cart is empty")
	CardTokenIsRequiredError    = New("The card token is required")
	ChangeShouldBePositiveError = New("The change should be positive")
)
