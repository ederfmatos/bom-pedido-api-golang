package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type MerchantMemoryRepository struct {
	merchants           map[string]*entity.Merchant
	merchantsByTenantId map[string]*entity.Merchant
}

func NewMerchantMemoryRepository() *MerchantMemoryRepository {
	return &MerchantMemoryRepository{
		merchants:           make(map[string]*entity.Merchant),
		merchantsByTenantId: make(map[string]*entity.Merchant),
	}
}

func (r *MerchantMemoryRepository) Create(_ context.Context, merchant *entity.Merchant) error {
	r.merchants[merchant.Id] = merchant
	r.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (r *MerchantMemoryRepository) Update(_ context.Context, merchant *entity.Merchant) error {
	r.merchants[merchant.Id] = merchant
	r.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (r *MerchantMemoryRepository) FindByTenantId(_ context.Context, tenantId string) (*entity.Merchant, error) {
	return r.merchantsByTenantId[tenantId], nil
}

func (r *MerchantMemoryRepository) IsActive(_ context.Context, merchantId string) (bool, error) {
	if merchant, exists := r.merchants[merchantId]; exists {
		return merchant.IsActive(), nil
	}
	return false, nil
}
