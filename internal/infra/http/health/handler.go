package health

import (
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/pkg/mongo"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Output struct {
	Ok bool `json:"ok"`
}

func Handle(client *redis.Client, mongoDatabase *mongo.Database) func(context echo.Context) error {
	return func(context echo.Context) error {
		healthErrorComposite := errors.NewCompositeError()
		ctx := context.Request().Context()
		if err := mongoDatabase.Ping(ctx); err != nil {
			healthErrorComposite.AppendError(err)
		}
		if err := client.Ping(ctx).Err(); err != nil {
			healthErrorComposite.AppendError(err)
		}
		if healthErrorComposite.HasError() {
			return healthErrorComposite.AsError()
		}
		return context.JSON(200, Output{Ok: true})
	}
}
