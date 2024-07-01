package entity

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"time"
)

var (
	ShoppingCartEmptyError      = errors.New("Your shopping cart is empty")
	CardTokenIsRequiredError    = errors.New("The card token is required")
	ChangeShouldBePositiveError = errors.New("The change should be positive")
)

type ShoppingCart struct {
	CustomerId string
	Items      []ShoppingCartItem
}

type ShoppingCartItem struct {
	Id        string
	ProductId string
	Quantity    int
	Observation string
	Price       float64
}

func (item *ShoppingCartItem) GetTotalPrice() float64 {
	return float64(item.Quantity) * item.Price
}

func NewShoppingCart(customerId string) *ShoppingCart {
	return &ShoppingCart{
		CustomerId: customerId,
		Items:      []ShoppingCartItem{},
	}
}

func (shoppingCart *ShoppingCart) AddItem(product *Product, quantity int, observation string) error {
	if product.IsUnAvailable() {
		return ProductUnAvailableError
	}
	item := ShoppingCartItem{
		Id:          value_object.NewID(),
		ProductId:   product.Id,
		Quantity:    quantity,
		Observation: observation,
		Price:       product.Price,
	}
	shoppingCart.Items = append(shoppingCart.Items, item)
	return nil
}

func (shoppingCart *ShoppingCart) GetPrice() float64 {
	totalPrice := float64(0)
	for _, item := range shoppingCart.Items {
		totalPrice += item.GetTotalPrice()
	}
	return totalPrice
}

func (shoppingCart *ShoppingCart) GetItems() []ShoppingCartItem {
	return shoppingCart.Items
}

func (shoppingCart *ShoppingCart) Checkout(
	paymentMethodString, deliveryModeString, paymentModeString, cardToken string,
	change float64,
	products map[string]*Product,
	deliveryTime time.Duration,
) (*Order, error) {
	if shoppingCart.IsEmpty() {
		return nil, ShoppingCartEmptyError
	}
	order, err := NewOrder(shoppingCart.CustomerId, paymentMethodString, paymentModeString, deliveryModeString, cardToken, change, time.Now().Add(deliveryTime))
	if err != nil {
		return nil, err
	}
	compositeError := errors.NewCompositeError()
	for _, item := range shoppingCart.Items {
		product := products[item.ProductId]
		err := order.AddProduct(product, item.Quantity, item.Observation)
		compositeError.Append(err)
	}
	if compositeError.HasError() {
		return nil, compositeError
	}
	return order, nil
}

func (shoppingCart *ShoppingCart) IsEmpty() bool {
	return len(shoppingCart.Items) == 0
}
