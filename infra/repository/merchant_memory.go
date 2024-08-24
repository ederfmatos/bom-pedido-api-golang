package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/merchant"
	"context"
)

type MerchantMemoryRepository struct {
	merchants           map[string]*merchant.Merchant
	merchantsByTenantId map[string]*merchant.Merchant
}

func NewMerchantMemoryRepository() repository.MerchantRepository {
	return &MerchantMemoryRepository{
		merchants:           make(map[string]*merchant.Merchant),
		merchantsByTenantId: make(map[string]*merchant.Merchant),
	}
}

func (repository *MerchantMemoryRepository) Create(_ context.Context, merchant *merchant.Merchant) error {
	repository.merchants[merchant.Id] = merchant
	repository.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (repository *MerchantMemoryRepository) Update(_ context.Context, merchant *merchant.Merchant) error {
	repository.merchants[merchant.Id] = merchant
	repository.merchantsByTenantId[merchant.TenantId] = merchant
	return nil
}

func (repository *MerchantMemoryRepository) FindByTenantId(_ context.Context, tenantId string) (*merchant.Merchant, error) {
	return repository.merchantsByTenantId[tenantId], nil
}

func (repository *MerchantMemoryRepository) IsActive(_ context.Context, merchantId string) (bool, error) {
	if aMerchant, exists := repository.merchants[merchantId]; exists {
		return aMerchant.IsActive(), nil
	}
	return false, nil
}
