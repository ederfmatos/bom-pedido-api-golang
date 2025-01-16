package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type CategoriesMongoRepository struct {
	*mongo.Collection
}

func NewCategoriesMongoRepository(database *mongo.Database) *CategoriesMongoRepository {
	return &CategoriesMongoRepository{Collection: database.ForCollection("product_categories")}
}

func (r *CategoriesMongoRepository) Create(ctx context.Context, category *entity.Category) error {
	return r.InsertOne(ctx, category)
}

func (r *CategoriesMongoRepository) ExistsById(ctx context.Context, id string) (bool, error) {
	return r.ExistsByID(ctx, id)
}

func (r *CategoriesMongoRepository) ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error) {
	return r.ExistsByTenantIdAnd(ctx, tenantId, "name", name)
}
