package shopping_cart

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"time"
)

type (
	ShoppingCart struct {
		CustomerId string             `bson:"_id"`
		TenantId   string             `bson:"tenantId"`
		Items      []ShoppingCartItem `bson:"items"`
	}
	ShoppingCartItem struct {
		Id          string  `bson:"id"`
		ProductId   string  `bson:"productId"`
		Quantity    int     `bson:"quantity"`
		Observation string  `bson:"observation"`
		Price       float64 `bson:"price"`
	}
)

func New(customerId, tenantId string) *ShoppingCart {
	return &ShoppingCart{
		CustomerId: customerId,
		TenantId:   tenantId,
		Items:      []ShoppingCartItem{},
	}
}

func (shoppingCart *ShoppingCart) AddItem(product *product.Product, quantity int, observation string) error {
	if product.IsUnAvailable() {
		return errors.ProductUnAvailableError
	}
	if quantity < 1 {
		return errors.New("quantity must be greater than 0")
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

func (item *ShoppingCartItem) GetTotalPrice() float64 {
	return float64(item.Quantity) * item.Price
}

func (shoppingCart *ShoppingCart) GetItems() []ShoppingCartItem {
	return shoppingCart.Items
}

func (shoppingCart *ShoppingCart) Checkout(
	paymentMethodString, deliveryModeString, paymentModeString, cardToken string,
	payback float64,
	products map[string]*product.Product,
	deliveryTime time.Duration,
) (*order.Order, error) {
	if shoppingCart.IsEmpty() {
		return nil, errors.ShoppingCartEmptyError
	}
	price := shoppingCart.GetPrice()
	anOrder, err := order.New(shoppingCart.CustomerId, paymentMethodString, paymentModeString, deliveryModeString, cardToken, payback, price, time.Now().Add(deliveryTime), shoppingCart.TenantId)
	if err != nil {
		return nil, err
	}
	compositeError := errors.NewCompositeError()
	for _, item := range shoppingCart.Items {
		aProduct := products[item.ProductId]
		err = anOrder.AddProduct(aProduct, item.Quantity, item.Observation)
		compositeError.Append(err)
	}
	if compositeError.HasError() {
		return nil, compositeError
	}
	return anOrder, nil
}

func (shoppingCart *ShoppingCart) IsEmpty() bool {
	return len(shoppingCart.Items) == 0
}

func CloneOrder(order *order.Order) *ShoppingCart {
	shoppingCart := &ShoppingCart{
		CustomerId: order.CustomerID,
		Items:      make([]ShoppingCartItem, len(order.Items)),
	}
	for i, item := range order.Items {
		shoppingCart.Items[i] = ShoppingCartItem{
			Id:          item.Id,
			ProductId:   item.ProductId,
			Quantity:    item.Quantity,
			Observation: item.Observation,
			Price:       item.Price,
		}
	}
	return shoppingCart
}
