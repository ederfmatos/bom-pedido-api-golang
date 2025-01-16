package entity

import (
	"bom-pedido-api/internal/domain/value_object"
)

type Category struct {
	Id          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	TenantId    string `bson:"tenantId"`
}

func NewCategory(name, description, tenantId string) *Category {
	return &Category{
		Id:          value_object.NewID(),
		Name:        name,
		Description: description,
		TenantId:    tenantId,
	}
}
