package messaging

import (
	"bom-pedido-api/application/event"
	"log/slog"
)

func HandleCreateProductProjection() func(event event.Event) error {
	return func(event event.Event) error {
		slog.Info("Handling CreateProductProjection event", "product", event.Data)
		return nil
	}
}
