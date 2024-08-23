package query

import (
	"bom-pedido-api/application/query"
	"bom-pedido-api/infra/repository"
	"context"
)

type ProductSqlQuery struct {
	repository.SqlConnection
}

func NewProductSqlQuery(sqlConnection repository.SqlConnection) query.ProductQuery {
	return &ProductSqlQuery{SqlConnection: sqlConnection}
}

func (sqlQuery *ProductSqlQuery) List(ctx context.Context, tenantId string) ([]query.Product, error) {
	products := make([]query.Product, 0)
	err := sqlQuery.Sql("select id, name, description, price from products where tenant_id = $1").
		Values(tenantId).
		List(ctx, func(getValues func(dest ...any) error) error {
			var product query.Product
			err := getValues(&product.Id, &product.Name, &product.Description, &product.Price)
			if err != nil {
				return err
			}
			products = append(products, product)
			return nil
		})
	return products, err
}
