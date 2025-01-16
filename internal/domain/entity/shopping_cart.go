package entity

import (
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"time"
)

type (
	ShoppingCart struct {
		CustomerId string                      `bson:"id"`
		TenantId   string                      `bson:"tenantId"`
		Items      map[string]ShoppingCartItem `bson:"items"`
	}
	ShoppingCartItem struct {
		Id          string  `bson:"id"`
		ProductId   string  `bson:"productId"`
		Quantity    int     `bson:"quantity"`
		Observation string  `bson:"observation"`
		Price       float64 `bson:"price"`
		ProductName string  `bson:"productName"`
	}
)

func NewShoppingCart(customerId, tenantId string) *ShoppingCart {
	return &ShoppingCart{
		CustomerId: customerId,
		TenantId:   tenantId,
		Items:      make(map[string]ShoppingCartItem),
	}
}

func NewShoppingCartFromOrder(order *Order) *ShoppingCart {
	shoppingCart := &ShoppingCart{
		CustomerId: order.CustomerID,
		Items:      make(map[string]ShoppingCartItem, len(order.Items)),
	}
	for _, item := range order.Items {
		shoppingCart.Items[item.Id] = ShoppingCartItem{
			Id:          item.Id,
			ProductId:   item.ProductId,
			Quantity:    item.Quantity,
			Observation: item.Observation,
			Price:       item.Price,
		}
	}
	return shoppingCart
}

func (shoppingCart *ShoppingCart) AddItem(product *Product, quantity int, observation string) error {
	if product.IsUnAvailable() {
		return errors.ProductUnAvailableError
	}
	if quantity < 1 {
		return errors.New("quantity must be greater than 0")
	}
	item := ShoppingCartItem{
		Id:          value_object.NewID(),
		ProductId:   product.Id,
		ProductName: product.Name,
		Quantity:    quantity,
		Observation: observation,
		Price:       product.Price,
	}
	if shoppingCart.Items == nil {
		shoppingCart.Items = make(map[string]ShoppingCartItem)
	}
	shoppingCart.Items[item.Id] = item
	return nil
}

func (shoppingCart *ShoppingCart) GetPrice() float64 {
	totalPrice := float64(0)
	for _, item := range shoppingCart.Items {
		totalPrice += item.GetTotalPrice()
	}
	return totalPrice
}

func (item *ShoppingCartItem) GetTotalPrice() float64 {
	return float64(item.Quantity) * item.Price
}

func (shoppingCart *ShoppingCart) Checkout(paymentMethodString, deliveryModeString, paymentModeString, cardToken string, payback float64, products map[string]*Product, deliveryTime time.Duration, merchantId string) (*Order, error) {
	if shoppingCart.IsEmpty() {
		return nil, errors.ShoppingCartEmptyError
	}
	price := shoppingCart.GetPrice()
	order, err := NewOrder(shoppingCart.CustomerId, paymentMethodString, paymentModeString, deliveryModeString, cardToken, payback, price, time.Now().Add(deliveryTime), merchantId)
	if err != nil {
		return nil, err
	}
	compositeError := errors.NewCompositeError()
	for _, item := range shoppingCart.Items {
		product := products[item.ProductId]
		err = order.AddProduct(product, item.Quantity, item.Observation)
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

func (shoppingCart *ShoppingCart) RemoveItem(id string) {
	delete(shoppingCart.Items, id)
}
