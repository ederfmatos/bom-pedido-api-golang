package repository

import (
	"bom-pedido-api/domain/entity/admin"
	"context"
)

type AdminRepository interface {
	FindByEmail(ctx context.Context, email string) (*admin.Admin, error)
	Create(ctx context.Context, admin *admin.Admin) error
}
