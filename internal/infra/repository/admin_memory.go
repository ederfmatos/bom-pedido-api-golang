package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type AdminMemoryRepository struct {
	admins map[string]*entity.Admin
}

func NewAdminMemoryRepository() *AdminMemoryRepository {
	return &AdminMemoryRepository{admins: make(map[string]*entity.Admin)}
}

func (repository *AdminMemoryRepository) Create(_ context.Context, admin *entity.Admin) error {
	repository.admins[admin.Id] = admin
	return nil
}

func (repository *AdminMemoryRepository) FindByEmail(_ context.Context, email string) (*entity.Admin, error) {
	for _, admin := range repository.admins {
		if admin.GetEmail() == email {
			return admin, nil
		}
	}
	return nil, nil
}
