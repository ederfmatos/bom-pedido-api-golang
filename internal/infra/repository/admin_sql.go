package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/admin"
	"context"
)

const (
	createAdminSQL      = "INSERT INTO admins (id, name, email, merchant_id) VALUES ($1, $2, $3, $4)"
	findAdminByEmailSQL = "SELECT id, name, email, merchant_id FROM admins WHERE email = $1"
)

type DefaultAdminRepository struct {
	SqlConnection
}

func NewDefaultAdminRepository(connection SqlConnection) repository.AdminRepository {
	return &DefaultAdminRepository{connection}
}

func (repository *DefaultAdminRepository) Create(ctx context.Context, admin *admin.Admin) error {
	return repository.Sql(createAdminSQL).
		Values(admin.Id, admin.Name, admin.GetEmail(), admin.MerchantId).
		Update(ctx)
}

func (repository *DefaultAdminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var id, name, merchantID string
	found, err := repository.Sql(findAdminByEmailSQL).
		Values(email).
		FindOne(ctx, &id, &name, &email, &merchantID)
	if err != nil || !found {
		return nil, err
	}
	return admin.Restore(id, name, email, merchantID)
}
