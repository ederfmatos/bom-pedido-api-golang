package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/admin"
	"context"
)

const (
	createAdminSql      = "INSERT INTO admins (id, name, email, merchant_id) VALUES ($1, $2, $3, $4)"
	findAdminByEmailSql = "SELECT id, name, email, merchant_id FROM admins WHERE email = $1"
)

type DefaultAdminRepository struct {
	SqlConnection
}

func NewDefaultAdminRepository(connection SqlConnection) repository.AdminRepository {
	return &DefaultAdminRepository{connection}
}

func (repository *DefaultAdminRepository) Create(ctx context.Context, admin *admin.Admin) error {
	return repository.Sql(createAdminSql).
		Values(admin.Id, admin.Name, admin.GetEmail(), admin.MerchantId).
		Update(ctx)
}

func (repository *DefaultAdminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var id, name, merchantId string
	found, err := repository.Sql(findAdminByEmailSql).
		Values(email).
		FindOne(ctx, &id, &name, &email, &merchantId)
	if err != nil || !found {
		return nil, err
	}
	return admin.Restore(id, name, email, merchantId)
}
