package query

import (
	"bom-pedido-api/internal/application/projection"
	"bom-pedido-api/internal/application/query"
	"bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
)

type ProductQuery struct {
	*mongo.Collection
}

func NewProductQuery(mongoDatabase *mongo.Database) query.ProductQuery {
	return &ProductQuery{Collection: mongoDatabase.ForCollection("products")}
}

func (q *ProductQuery) List(ctx context.Context, filter projection.ProductListFilter) (*projection.Page[projection.ProductListItem], error) {
	mongoFilter := map[string]interface{}{"tenantId": filter.TenantId}

	totalItems, err := q.Collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("count products: %v", err)
	}

	skip := (filter.CurrentPage - 1) * filter.PageSize
	limit := filter.PageSize

	cursor, err := q.Collection.Find(ctx, mongoFilter, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("find products: %v", err)
	}
	defer cursor.Close(ctx)

	totalPages := (totalItems + filter.PageSize - 1) / filter.PageSize

	size := filter.CurrentPage
	if totalItems < size {
		size = totalItems
	}

	page := projection.Page[projection.ProductListItem]{
		CurrentPage: filter.CurrentPage,
		PageSize:    filter.PageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		LastPage:    totalPages == filter.CurrentPage,
		Items:       make([]projection.ProductListItem, size),
	}

	if err = cursor.All(ctx, &page.Items); err != nil {
		return nil, fmt.Errorf("decode products: %v", err)
	}

	return &page, nil
}
