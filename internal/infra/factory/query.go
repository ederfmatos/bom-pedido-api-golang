package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/query"
	"bom-pedido-api/internal/infra/repository"
)

func queryFactory(connection repository.SqlConnection) *factory.QueryFactory {
	return factory.NewQueryFactory(query.NewProductSqlQuery(connection))
}
