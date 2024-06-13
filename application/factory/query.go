package factory

import (
	"bom-pedido-api/application/query"
)

type QueryFactory struct {
	ProductQuery query.ProductQuery
}

func NewQueryFactory(
	productQuery query.ProductQuery,
) *QueryFactory {
	return &QueryFactory{
		ProductQuery: productQuery,
	}
}
