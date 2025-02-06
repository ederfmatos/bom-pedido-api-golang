package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type AdminMongoRepository struct {
	collection mongo.Collection
}

func NewAdminMongoRepository(database *mongo.Database) *AdminMongoRepository {
	return &AdminMongoRepository{collection: database.ForCollection("admins")}
}

func (r *AdminMongoRepository) Create(ctx context.Context, admin *entity.Admin) error {
	return r.collection.InsertOne(ctx, admin)
}

func (r *AdminMongoRepository) FindByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.collection.FindBy(ctx, "email", email, &admin)
	if err != nil || admin.Id == "" {
		return nil, err
	}
	return &admin, nil
}
