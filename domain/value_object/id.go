package value_object

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func NewID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func NewTenantId() string {
	return ulid.Make().String()
}
