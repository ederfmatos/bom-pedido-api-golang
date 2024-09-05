package query

import (
	"bom-pedido-api/application/projection"
	"context"
)

type ProductQuery interface {
	List(ctx context.Context, filter projection.ProductListFilter) (*projection.Page[projection.ProductListItem], error)
}
