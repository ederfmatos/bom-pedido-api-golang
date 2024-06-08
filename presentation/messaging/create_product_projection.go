package messaging

import (
	"bom-pedido-api/application/event"
	"fmt"
)

func HandleCreateProductProjection() func(event event.Event) error {
	return func(event event.Event) error {
		fmt.Printf("Event received: %v\n", event)
		return nil
	}
}
