package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/merchant"
	"context"
)

const (
	sqlCreateMerchant         = "INSERT INTO merchants (id, name, email, phone_number, tenant_id, domain, status) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	sqlUpdateMerchant         = "UPDATE merchants SET name = $1, email = $2, phone_number = $3, status = $4, domain = $5 WHERE id = $6;"
	sqlMerchantIdActive       = "SELECT 1 FROM merchants WHERE id = $1 AND status = $2 LIMIT 1;"
	sqlFindMerchantByTenantId = "SELECT id, name, email, phone_number, domain, status FROM merchants WHERE tenant_id = $1 LIMIT 1;"
)

type DefaultMerchantRepository struct {
	SqlConnection
}

func NewDefaultMerchantRepository(sqlConnection SqlConnection) repository.MerchantRepository {
	return &DefaultMerchantRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultMerchantRepository) Create(ctx context.Context, merchant *merchant.Merchant) error {
	return repository.Sql(sqlCreateMerchant).
		Values(merchant.Id, merchant.Name, merchant.Email.Value(), merchant.PhoneNumber.Value(), merchant.TenantId, merchant.Domain, merchant.Status).
		Update(ctx)
}

func (repository *DefaultMerchantRepository) Update(ctx context.Context, merchant *merchant.Merchant) error {
	return repository.Sql(sqlUpdateMerchant).
		Values(merchant.Name, merchant.Email.Value(), merchant.PhoneNumber.Value(), merchant.Status, merchant.Domain, merchant.Id).
		Update(ctx)
}

func (repository *DefaultMerchantRepository) FindByTenantId(ctx context.Context, tenantId string) (*merchant.Merchant, error) {
	var id, name, email, phoneNumber, domain, status string
	found, err := repository.Sql(sqlFindMerchantByTenantId).
		Values(tenantId).
		FindOne(ctx, &id, &name, &email, &phoneNumber, &domain, &status)
	if err != nil || !found {
		return nil, err
	}
	return merchant.Restore(id, name, email, phoneNumber, status, domain, tenantId)
}

func (repository *DefaultMerchantRepository) IsActive(ctx context.Context, merchantId string) (bool, error) {
	return repository.Sql(sqlMerchantIdActive).Values(merchantId, merchant.ACTIVE).Exists(ctx)
}
