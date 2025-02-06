package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type MerchantMongoRepository struct {
	collection mongo.Collection
}

func NewMerchantMongoRepository(database *mongo.Database) *MerchantMongoRepository {
	return &MerchantMongoRepository{collection: database.ForCollection("merchants")}
}

func (r *MerchantMongoRepository) Create(ctx context.Context, merchant *entity.Merchant) error {
	return r.collection.InsertOne(ctx, merchant)
}

func (r *MerchantMongoRepository) Update(ctx context.Context, merchant *entity.Merchant) error {
	return r.collection.UpdateByID(ctx, merchant.Id, merchant)
}

func (r *MerchantMongoRepository) FindByTenantId(ctx context.Context, tenantId string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.collection.FindBy(ctx, "tenantId", tenantId, &merchant)
	if err != nil || merchant.Id == "" {
		return nil, err
	}
	return &merchant, nil
}

func (r *MerchantMongoRepository) IsActive(ctx context.Context, merchantId string) (bool, error) {
	var merchant entity.Merchant
	err := r.collection.FindByID(ctx, merchantId, &merchant)
	if err != nil || merchant.Id == "" {
		return false, err
	}
	return merchant.IsActive(), nil
}
