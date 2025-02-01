package category

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"context"
	"fmt"
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

	CreateCategoryOutput struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func NewCreateCategory(factory *factory.ApplicationFactory) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepository: factory.ProductCategoryRepository,
	}
}

func (u *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (*CreateCategoryOutput, error) {
	existsByName, err := u.categoryRepository.ExistsByNameAndTenantId(ctx, input.Name, input.TenantId)
	if err != nil {
		return nil, err
	}
	if existsByName {
		return nil, CategoryWithSameNameError
	}

	category := entity.NewCategory(input.Name, input.Description, input.TenantId)
	if err = u.categoryRepository.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("create category: %v", err)
	}

	return &CreateCategoryOutput{
		ID:          category.Id,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}
