package query

import (
	"bom-pedido-api/application/projection"
	"bom-pedido-api/application/query"
	"bom-pedido-api/infra/repository"
	"context"
	"math"
)

type ProductSqlQuery struct {
	repository.SqlConnection
}

func NewProductSqlQuery(sqlConnection repository.SqlConnection) query.ProductQuery {
	return &ProductSqlQuery{SqlConnection: sqlConnection}
}

func (q *ProductSqlQuery) List(ctx context.Context, filter projection.ProductListFilter) (*projection.Page[projection.ProductListItem], error) {
	page := projection.Page[projection.ProductListItem]{
		CurrentPage: filter.CurrentPage,
		PageSize:    filter.PageSize,
		TotalPages:  0,
		TotalItems:  0,
		LastPage:    false,
		Items:       make([]projection.ProductListItem, 0),
	}
	err := q.Sql("SELECT COUNT(id) as TOTAL_ITEMS FROM products WHERE tenant_id = $1").Values(filter.TenantId).Count(ctx, &page.TotalItems)
	if err != nil || page.TotalItems == 0 {
		return nil, err
	}
	page.TotalPages = int32(math.Ceil(float64(page.TotalItems) / float64(filter.PageSize)))
	page.LastPage = filter.CurrentPage == page.TotalPages
	skip := calculateSkip(filter)
	err = q.Sql("select id, name, description, price, status from products where tenant_id = $1 LIMIT $2 OFFSET $3").
		Values(filter.TenantId, filter.PageSize, skip).
		List(ctx, func(getValues func(dest ...any) error) error {
			var product projection.ProductListItem
			err = getValues(&product.Id, &product.Name, &product.Description, &product.Price, &product.Status)
			if err == nil {
				page.Items = append(page.Items, product)
			}
			return err
		})
	return &page, err
}

func calculateSkip(filter projection.ProductListFilter) int32 {
	if filter.CurrentPage > 1 {
		return (filter.CurrentPage - 1) * filter.PageSize
	}
	return 0
}
