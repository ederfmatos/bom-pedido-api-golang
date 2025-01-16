package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type AdminRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.Admin, error)
	Create(ctx context.Context, admin *entity.Admin) error
}
