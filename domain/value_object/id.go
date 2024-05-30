package value_object

import "github.com/google/uuid"

func NewID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
