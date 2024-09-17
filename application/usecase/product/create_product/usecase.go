package create_product

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/errors"
	"context"
)

type (
	UseCase struct {
		productRepository  repository.ProductRepository
		categoryRepository repository.ProductCategoryRepository
		eventEmitter       event.Emitter
	}
	Input struct {
		Name        string
		Description string
		CategoryId  string
		TenantId    string
		Price       float64
	}
	Output struct {
		Id string `json:"id"`
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		productRepository:  factory.ProductRepository,
		categoryRepository: factory.ProductCategoryRepository,
		eventEmitter:       factory.EventEmitter,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
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
	aProduct, err := product.New(input.Name, input.Description, input.Price, input.CategoryId, input.TenantId)
	if err != nil {
		return nil, err
	}
	err = useCase.productRepository.Create(ctx, aProduct)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(ctx, event.NewProductCreatedEvent(aProduct))
	if err != nil {
		return nil, err
	}
	return &Output{Id: aProduct.Id}, nil
}
