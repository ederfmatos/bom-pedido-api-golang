package repository

import (
	"bom-pedido-api/internal/domain/entity/admin"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type AdminMongoRepository struct {
	collection *mongo.Collection
}

func NewAdminMongoRepository(database *mongo.Database) *AdminMongoRepository {
	return &AdminMongoRepository{collection: database.ForCollection("admins")}
}

func (r *AdminMongoRepository) Create(ctx context.Context, admin *admin.Admin) error {
	return r.collection.InsertOne(ctx, admin)
}

func (r *AdminMongoRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var anAdmin admin.Admin
	err := r.collection.FindBy(ctx, "email", email, &anAdmin)
	if err != nil || anAdmin.Id == "" {
		return nil, err
	}
	return &anAdmin, nil
}
