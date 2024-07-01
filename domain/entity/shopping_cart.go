package entity

import (
	"bom-pedido-api/domain/value_object"
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
