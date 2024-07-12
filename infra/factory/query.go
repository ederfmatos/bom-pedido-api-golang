package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/query"
	"bom-pedido-api/infra/repository"
)

func queryFactory(connection repository.SqlConnection) *factory.QueryFactory {
	return factory.NewQueryFactory(query.NewProductSqlQuery(connection))
}
