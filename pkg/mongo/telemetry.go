package mongo

import (
	"bom-pedido-api/pkg/telemetry"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TelemetryCollection struct {
	collection Collection
}

func NewTelemetryCollection(collection Collection) *TelemetryCollection {
	return &TelemetryCollection{collection: collection}
}

func (c TelemetryCollection) DeleteByID(ctx context.Context, id string) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.DeleteByID", func(ctx context.Context) error {
		return c.collection.DeleteByID(ctx, id)
	}, "record.id", id)
}

func (c TelemetryCollection) Upsert(ctx context.Context, id string, value any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.Upsert", func(ctx context.Context) error {
		return c.collection.Upsert(ctx, id, value)
	}, "record.id", id)
}

func (c TelemetryCollection) UpdateByID(ctx context.Context, id string, value any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.UpdateByID", func(ctx context.Context) error {
		return c.collection.UpdateByID(ctx, id, value)
	}, "record.id", id)
}

func (c TelemetryCollection) FindByID(ctx context.Context, id string, target any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.FindByID", func(ctx context.Context) error {
		return c.collection.FindByID(ctx, id, target)
	}, "record.id", id)
}

func (c TelemetryCollection) FindByTenantIdAnd(ctx context.Context, tenantId, param, value string, target any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.FindByTenantIdAnd", func(ctx context.Context) error {
		return c.collection.FindByTenantIdAnd(ctx, tenantId, param, value, target)
	})
}

func (c TelemetryCollection) FindBy(ctx context.Context, param, value string, target any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.FindBy", func(ctx context.Context) error {
		return c.collection.FindBy(ctx, param, value, target)
	})
}

func (c TelemetryCollection) FindByValues(ctx context.Context, values map[string]interface{}, target any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.FindByValues", func(ctx context.Context) error {
		return c.collection.FindByValues(ctx, values, target)
	})
}

func (c TelemetryCollection) FindAllByID(ctx context.Context, ids []string) (*mongo.Cursor, error) {
	return telemetry.StartSpan[*mongo.Cursor](ctx, "Mongo.FindAllByID", func(ctx context.Context) (*mongo.Cursor, error) {
		return c.collection.FindAllByID(ctx, ids)
	})
}

func (c TelemetryCollection) FindAllBy(ctx context.Context, values map[string]interface{}) (*mongo.Cursor, error) {
	return telemetry.StartSpan[*mongo.Cursor](ctx, "Mongo.FindAllBy", func(ctx context.Context) (*mongo.Cursor, error) {
		return c.collection.FindAllBy(ctx, values)
	})
}

func (c TelemetryCollection) Find(ctx context.Context, filter map[string]interface{}, skip, limit int64) (*mongo.Cursor, error) {
	return telemetry.StartSpan[*mongo.Cursor](ctx, "Mongo.Find", func(ctx context.Context) (*mongo.Cursor, error) {
		return c.collection.Find(ctx, filter, skip, limit)
	})
}

func (c TelemetryCollection) InsertOne(ctx context.Context, value any) error {
	return telemetry.StartSpanReturningError(ctx, "Mongo.InsertOne", func(ctx context.Context) error {
		return c.collection.InsertOne(ctx, value)
	})
}

func (c TelemetryCollection) ExistsByID(ctx context.Context, id string) (bool, error) {
	return telemetry.StartSpan[bool](ctx, "Mongo.ExistsByID", func(ctx context.Context) (bool, error) {
		return c.collection.ExistsByID(ctx, id)
	}, "record.id", id)
}

func (c TelemetryCollection) ExistsBy(ctx context.Context, name string, value string) (bool, error) {
	return telemetry.StartSpan[bool](ctx, "Mongo.ExistsBy", func(ctx context.Context) (bool, error) {
		return c.collection.ExistsBy(ctx, name, value)
	})
}

func (c TelemetryCollection) ExistsByTenantIdAnd(ctx context.Context, tenantId, name, value string) (bool, error) {
	return telemetry.StartSpan[bool](ctx, "Mongo.ExistsByTenantIdAnd", func(ctx context.Context) (bool, error) {
		return c.collection.ExistsByTenantIdAnd(ctx, tenantId, name, value)
	})
}

func (c TelemetryCollection) CountDocuments(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return telemetry.StartSpan[int64](ctx, "Mongo.CountDocuments", func(ctx context.Context) (int64, error) {
		return c.collection.CountDocuments(ctx, filter)
	})
}

func (c TelemetryCollection) FetchStream(ctx context.Context) (<-chan string, error) {
	return c.collection.FetchStream(ctx)
}
