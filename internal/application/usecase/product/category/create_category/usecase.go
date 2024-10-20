package create_category

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

var (
	CategoryWithSameNameError = errors.New("Already exists a category with the same name")
)

type (
	UseCase struct {
		categoryRepository repository.ProductCategoryRepository
	}

	Input struct {
		TenantId    string
		Name        string
		Description string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		categoryRepository: factory.ProductCategoryRepository,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) error {
	existsByName, err := u.categoryRepository.ExistsByNameAndTenantId(ctx, input.Name, input.TenantId)
	if err != nil {
		return err
	}
	if existsByName {
		return CategoryWithSameNameError
	}
	category := product.NewCategory(input.Name, input.Description, input.TenantId)
	return u.categoryRepository.Create(ctx, category)
}
