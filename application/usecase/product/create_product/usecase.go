package create_product

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/events"
	"context"
)

type (
	UseCase struct {
		productRepository repository.ProductRepository
		eventEmitter      event.Emitter
	}
	Input struct {
		Context     context.Context
		Name        string
		Description string
		Price       float64
	}
	Output struct {
		Id string `json:"id"`
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		productRepository: factory.ProductRepository,
		eventEmitter:      factory.EventEmitter,
	}
}

func (useCase *UseCase) Execute(input Input) (*Output, error) {
	existsByName, err := useCase.productRepository.ExistsByName(input.Context, input.Name)
	if err != nil {
		return nil, err
	}
	if existsByName {
		return nil, errors.ProductWithSameNameError
	}
	product, err := product.New(input.Name, input.Description, input.Price)
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
	return &Output{Id: product.Id}, nil
}
