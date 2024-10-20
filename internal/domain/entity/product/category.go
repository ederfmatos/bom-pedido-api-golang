package product

import (
	"bom-pedido-api/internal/domain/value_object"
)

type Category struct {
	Id          string
	Name        string
	Description string
	TenantId    string
}

func NewCategory(name, description, tenantId string) *Category {
	return &Category{
		Id:          value_object.NewID(),
		Name:        name,
		Description: description,
		TenantId:    tenantId,
	}
}
