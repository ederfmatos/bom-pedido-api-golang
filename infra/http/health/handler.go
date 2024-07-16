package health

import (
	"bom-pedido-api/domain/errors"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Output struct {
	Ok bool `json:"ok"`
}

func Handle(database *sql.DB, client *redis.Client, mongoClient *mongo.Client) func(context echo.Context) error {
	return func(context echo.Context) error {
		healthErrorComposite := errors.NewCompositeError()
		ctx := context.Request().Context()
		if err := mongoClient.Ping(ctx, nil); err != nil {
			healthErrorComposite.Append(err)
		}
		if err := database.PingContext(ctx); err != nil {
			healthErrorComposite.Append(err)
		}
		if err := client.Ping(ctx).Err(); err != nil {
			healthErrorComposite.Append(err)
		}
		if healthErrorComposite.HasError() {
			return healthErrorComposite.AsError()
		}
		return context.JSON(200, Output{Ok: true})
	}
}
