package projection

type (
	ShoppingCartItem struct {
		Id          string  `json:"id"`
		ProductId   string  `json:"productId"`
		ProductName string  `json:"productName"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		TotalPrice  float64 `json:"totalPrice"`
	}

	ShoppingCart struct {
		Amount float64            `json:"amount"`
		Items  []ShoppingCartItem `json:"items"`
	}
)
