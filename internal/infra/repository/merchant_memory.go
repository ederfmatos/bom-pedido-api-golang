package repository

import (
	"bom-pedido-api/internal/domain/entity/merchant"
	"context"
)

type MerchantMemoryRepository struct {
	merchants           map[string]*merchant.Merchant
	merchantsByTenantId map[string]*merchant.Merchant
}

func NewMerchantMemoryRepository() *MerchantMemoryRepository {
	return &MerchantMemoryRepository{
		merchants:           make(map[string]*merchant.Merchant),
		merchantsByTenantId: make(map[string]*merchant.Merchant),
	}
}

func (r *MerchantMemoryRepository) Create(_ context.Context, merchant *merchant.Merchant) error {
	r.merchants[merchant.Id] = merchant
	r.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (r *MerchantMemoryRepository) Update(_ context.Context, merchant *merchant.Merchant) error {
	r.merchants[merchant.Id] = merchant
	r.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (r *MerchantMemoryRepository) FindByTenantId(_ context.Context, tenantId string) (*merchant.Merchant, error) {
	return r.merchantsByTenantId[tenantId], nil
}

func (r *MerchantMemoryRepository) IsActive(_ context.Context, merchantId string) (bool, error) {
	if aMerchant, exists := r.merchants[merchantId]; exists {
		return aMerchant.IsActive(), nil
	}
	return false, nil
}
