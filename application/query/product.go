package query

import "context"

type Product struct {
	Id          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	ImageURL    string  `json:"imageURL,omitempty"`
}

type ProductQuery interface {
	List(ctx context.Context) ([]Product, error)
}
