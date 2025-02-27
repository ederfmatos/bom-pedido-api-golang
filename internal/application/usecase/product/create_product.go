package product

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	CreateProductUseCase struct {
		productRepository  repository.ProductRepository
		categoryRepository repository.ProductCategoryRepository
		eventEmitter       event.Emitter
	}
	CreateProductInput struct {
		Name        string
		Description string
		CategoryId  string
		TenantId    string
		Price       float64
	}
	CreateProductOutput struct {
		Id string `json:"id"`
	}
)

func NewCreateProduct(factory *factory.ApplicationFactory) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepository:  factory.ProductRepository,
		categoryRepository: factory.ProductCategoryRepository,
		eventEmitter:       factory.EventEmitter,
	}
}

func (useCase *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*CreateProductOutput, error) {
	existsByName, err := useCase.productRepository.ExistsByNameAndTenantId(ctx, input.Name, input.TenantId)
	if err != nil {
		return nil, err
	}
	if existsByName {
		return nil, errors.ProductWithSameNameError
	}
	existsCategory, err := useCase.categoryRepository.ExistsById(ctx, input.CategoryId)
	if err != nil {
		return nil, err
	}
	if !existsCategory {
		return nil, errors.ProductCategoryNotFoundError
	}
	product, err := entity.NewProduct(input.Name, input.Description, input.Price, input.CategoryId, input.TenantId)
	if err != nil {
		return nil, err
	}
	err = useCase.productRepository.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(ctx, event.NewProductCreatedEvent(product))
	if err != nil {
		return nil, err
	}
	return &CreateProductOutput{Id: product.Id}, nil
}
