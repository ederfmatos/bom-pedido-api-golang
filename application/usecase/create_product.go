package usecase

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/events"
	"context"
)

type CreateProductInput struct {
	Context     context.Context
	Name        string
	Description string
	Price       float64
}

type CreateProductOutput struct {
	ID string `json:"id"`
}

type CreateProductUseCase struct {
	productRepository repository.ProductRepository
	eventEmitter      event.EventEmitter
}

func NewCreateProductUseCase(factory *factory.ApplicationFactory) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepository: factory.ProductRepository,
		eventEmitter:      factory.EventEmitter,
	}
}

func (useCase CreateProductUseCase) Execute(input CreateProductInput) (*CreateProductOutput, error) {
	existsByName, err := useCase.productRepository.ExistsByName(input.Context, input.Name)
	if err != nil {
		return nil, err
	}
	if existsByName {
		return nil, errors.ProductWithSameName
	}
	product, err := entity.NewProduct(input.Name, input.Description, input.Price)
	if err != nil {
		return nil, err
	}
	err = useCase.productRepository.Create(input.Context, product)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(input.Context, events.NewProductCreatedEvent(product))
	if err != nil {
		return nil, err
	}
	return &CreateProductOutput{ID: product.ID}, nil
}
