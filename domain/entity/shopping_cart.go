package entity

import (
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
)

var (
	ShoppingCartNotFoundError      = errors.New("Produto não encontrado")
	CreditCardTokenIsRequiredError = errors.New("The credit card token is required")
)

type ShoppingCart struct {
	CustomerId string
	Items      []ShoppingCartItem
}

type ShoppingCartItem struct {
	Id          string
	ProductId   string
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
		ProductId:   product.ID,
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

func (shoppingCart *ShoppingCart) Checkout(paymentMethod enums.PaymentMethod, deliveryMode enums.DeliveryMode, paymentMode enums.PaymentMode, change float64, cardToken string, products map[string]*Product) (*Order, error) {
	composite := errors.NewCompositeError()
	if paymentMethod.IsCreditCard() && paymentMode.IsValid() && cardToken == "" {
		composite.Append(CreditCardTokenIsRequiredError)
	}
	order := NewOrder(shoppingCart.CustomerId, paymentMethod, paymentMode, deliveryMode, cardToken, change)
	for _, item := range shoppingCart.Items {
		product := products[item.ProductId]
		err := order.AddProduct(product, item.Quantity, item.Observation)
		composite.Append(err)
	}
	if composite.HasError() {
		return nil, composite.AsError()
	}
	return order, nil
}
