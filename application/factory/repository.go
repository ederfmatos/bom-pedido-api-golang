package factory

import (
	"bom-pedido-api/application/repository"
)

type RepositoryFactory struct {
	CustomerRepository     repository.CustomerRepository
	ProductRepository      repository.ProductRepository
	ShoppingCartRepository repository.ShoppingCartRepository
}

func NewRepositoryFactory(
	customerRepository repository.CustomerRepository,
	productRepository repository.ProductRepository,
	shoppingCartRepository repository.ShoppingCartRepository,
) *RepositoryFactory {
	return &RepositoryFactory{
		CustomerRepository:     customerRepository,
		ProductRepository:      productRepository,
		ShoppingCartRepository: shoppingCartRepository,
	}
}
