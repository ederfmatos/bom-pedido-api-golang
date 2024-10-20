package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/admin"
	"context"
)

type AdminMemoryRepository struct {
	admins map[string]*admin.Admin
}

func NewAdminMemoryRepository() repository.AdminRepository {
	return &AdminMemoryRepository{admins: make(map[string]*admin.Admin)}
}

func (repository *AdminMemoryRepository) Create(_ context.Context, admin *admin.Admin) error {
	repository.admins[admin.Id] = admin
	return nil
}

func (repository *AdminMemoryRepository) FindByEmail(_ context.Context, email string) (*admin.Admin, error) {
	for _, anAdmin := range repository.admins {
		if anAdmin.GetEmail() == email {
			return anAdmin, nil
		}
	}
	return nil, nil
}
