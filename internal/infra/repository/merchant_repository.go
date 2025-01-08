package repository

import (
	"bom-pedido-api/internal/domain/entity/merchant"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type MerchantMongoRepository struct {
	collection *mongo.Collection
}

func NewMerchantMongoRepository(database *mongo.Database) *MerchantMongoRepository {
	return &MerchantMongoRepository{collection: database.ForCollection("merchants")}
}

func (r *MerchantMongoRepository) Create(ctx context.Context, merchant *merchant.Merchant) error {
	return r.collection.InsertOne(ctx, merchant)
}

func (r *MerchantMongoRepository) Update(ctx context.Context, merchant *merchant.Merchant) error {
	return r.collection.UpdateByID(ctx, merchant.Id, merchant)
}

func (r *MerchantMongoRepository) FindByTenantId(ctx context.Context, tenantId string) (*merchant.Merchant, error) {
	var aMerchant merchant.Merchant
	err := r.collection.FindBy(ctx, "tenantId", tenantId, &aMerchant)
	if err != nil || aMerchant.Id == "" {
		return nil, err
	}
	return &aMerchant, nil
}

func (r *MerchantMongoRepository) IsActive(ctx context.Context, merchantId string) (bool, error) {
	var aMerchant merchant.Merchant
	err := r.collection.FindByID(ctx, merchantId, &aMerchant)
	if err != nil || aMerchant.Id == "" {
		return false, err
	}
	return aMerchant.IsActive(), nil
}
