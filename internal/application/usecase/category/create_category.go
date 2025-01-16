package category

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

var (
	CategoryWithSameNameError = errors.New("Already exists a category with the same name")
)

type (
	CreateCategoryUseCase struct {
		categoryRepository repository.ProductCategoryRepository
	}

	CreateCategoryInput struct {
		TenantId    string
		Name        string
		Description string
	}
)

func NewCreateCategory(factory *factory.ApplicationFactory) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepository: factory.ProductCategoryRepository,
	}
}

func (u *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) error {
	existsByName, err := u.categoryRepository.ExistsByNameAndTenantId(ctx, input.Name, input.TenantId)
	if err != nil {
		return err
	}
	if existsByName {
		return CategoryWithSameNameError
	}
	category := entity.NewCategory(input.Name, input.Description, input.TenantId)
	return u.categoryRepository.Create(ctx, category)
}
