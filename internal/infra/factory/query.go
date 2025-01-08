package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/query"
	"bom-pedido-api/pkg/mongo"
)

func queryFactory(mongoDatabase *mongo.Database) *factory.QueryFactory {
	return factory.NewQueryFactory(query.NewProductQuery(mongoDatabase))
}
