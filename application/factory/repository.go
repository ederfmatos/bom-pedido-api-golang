package factory

import (
	"bom-pedido-api/application/repository"
)

type RepositoryFactory struct {
	CustomerRepository repository.CustomerRepository
	ProductRepository  repository.ProductRepository
}

func NewRepositoryFactory(
	customerRepository repository.CustomerRepository,
	productRepository repository.ProductRepository,
) *RepositoryFactory {
	return &RepositoryFactory{
		CustomerRepository: customerRepository,
		ProductRepository:  productRepository,
	}
}
