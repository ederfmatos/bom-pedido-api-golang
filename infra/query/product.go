package query

import (
	"bom-pedido-api/application/query"
	"bom-pedido-api/infra/repository"
	"context"
)

type ProductSqlQuery struct {
	repository.SqlConnection
}

func NewProductSqlQuery(sqlConnection repository.SqlConnection) *ProductSqlQuery {
	return &ProductSqlQuery{SqlConnection: sqlConnection}
}

func (sqlQuery *ProductSqlQuery) List(context context.Context) ([]query.Product, error) {
	var products []query.Product
	err := sqlQuery.Sql("select id, name, description, price from products").
		List(context, func(getValues func(dest ...any) error) error {
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
