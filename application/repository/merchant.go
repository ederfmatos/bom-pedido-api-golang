package repository

import (
	"bom-pedido-api/domain/entity/merchant"
	"context"
)

type MerchantRepository interface {
	FindByTenantId(ctx context.Context, tenantId string) (*merchant.Merchant, error)
	IsActive(ctx context.Context, merchantId string) (bool, error)
	Create(ctx context.Context, merchant *merchant.Merchant) error
	Update(ctx context.Context, merchant *merchant.Merchant) error
}
