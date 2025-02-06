package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Collection interface {
		DeleteByID(ctx context.Context, id string) error
		Upsert(ctx context.Context, id string, value any) error
		UpdateByID(ctx context.Context, id string, value any) error
		FindByID(ctx context.Context, id string, target any) error
		FindByTenantIdAnd(ctx context.Context, tenantId, param, value string, target any) error
		FindBy(ctx context.Context, param, value string, target any) error
		FindByValues(ctx context.Context, values map[string]interface{}, target any) error
		FindAllByID(ctx context.Context, ids []string) (*mongo.Cursor, error)
		FindAllBy(ctx context.Context, values map[string]interface{}) (*mongo.Cursor, error)
		Find(ctx context.Context, filter map[string]interface{}, skip, limit int64) (*mongo.Cursor, error)
		InsertOne(ctx context.Context, value any) error
		ExistsByID(ctx context.Context, id string) (bool, error)
		ExistsBy(ctx context.Context, name string, value string) (bool, error)
		FetchStream(ctx context.Context) (<-chan string, error)
		CountDocuments(ctx context.Context, filter map[string]interface{}) (int64, error)
		ExistsByTenantIdAnd(ctx context.Context, tenantId string, key string, value string) (bool, error)
	}

	collection struct {
		*mongo.Collection
	}
)

type Database struct {
	database *mongo.Database
}

func NewDatabase(database *mongo.Database) *Database {
	return &Database{database: database}
}

func (d Database) ForCollection(name string) Collection {
	return newCollection(d.database.Collection(name))
}

func (d Database) Ping(ctx context.Context) error {
	return d.database.Client().Ping(ctx, nil)
}

func (d Database) Disconnect(ctx context.Context) error {
	return d.database.Client().Disconnect(ctx)
}

func newCollection(mongoCollection *mongo.Collection) Collection {
	return NewTelemetryCollection(&collection{Collection: mongoCollection})
}

func (c collection) findOne(ctx context.Context, filter bson.M, target any) error {
	result := c.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return fmt.Errorf("find record: %w", err)
	}
	if err := result.Decode(target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

func (c collection) DeleteByID(ctx context.Context, id string) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

func (c collection) Upsert(ctx context.Context, id string, value any) error {
	return c.update(ctx, id, value, true)
}

func (c collection) UpdateByID(ctx context.Context, id string, value any) error {
	return c.update(ctx, id, value, false)
}

func (c collection) update(ctx context.Context, id string, value any, upsert bool) error {
	update := bson.M{"$set": value}
	updateOptions := options.Update().SetUpsert(upsert)
	filter := bson.D{{Key: "id", Value: id}}
	_, err := c.Collection.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return fmt.Errorf("update record: %w", err)
	}
	return nil
}

func (c collection) FindByID(ctx context.Context, id string, target any) error {
	return c.findOne(ctx, bson.M{"id": id}, target)
}

func (c collection) FindByTenantIdAnd(ctx context.Context, tenantId, param, value string, target any) error {
	return c.findOne(ctx, bson.M{"tenantId": tenantId, param: value}, target)
}

func (c collection) FindBy(ctx context.Context, param, value string, target any) error {
	return c.findOne(ctx, bson.M{param: value}, target)
}

func (c collection) FindByValues(ctx context.Context, values map[string]interface{}, target any) error {
	return c.findOne(ctx, values, target)
}

func (c collection) FindAllByID(ctx context.Context, ids []string) (*mongo.Cursor, error) {
	filter := bson.M{"id": bson.M{"$in": ids}}
	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find records: %w", err)
	}
	return cursor, nil
}

func (c collection) FindAllBy(ctx context.Context, values map[string]interface{}) (*mongo.Cursor, error) {
	filter := bson.M(values)
	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find records: %w", err)
	}
	return cursor, nil
}

func (c collection) Find(ctx context.Context, filter map[string]interface{}, skip, limit int64) (*mongo.Cursor, error) {
	cursor, err := c.Collection.Find(ctx, bson.M(filter), &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("find records: %w", err)
	}
	return cursor, nil
}

func (c collection) InsertOne(ctx context.Context, value any) error {
	_, err := c.Collection.InsertOne(ctx, value)
	if err != nil {
		return fmt.Errorf("insert record: %w", err)
	}
	return nil
}

func (c collection) ExistsByID(ctx context.Context, id string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{"id": id})
}

func (c collection) ExistsBy(ctx context.Context, name string, value string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{name: value})
}

func (c collection) existsByValues(ctx context.Context, values map[string]interface{}) (bool, error) {
	result := c.Collection.FindOne(ctx, bson.M(values))
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("find record: %w", err)
	}
	return true, nil
}

func (c collection) ExistsByTenantIdAnd(ctx context.Context, tenantId, name, value string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{"tenantId": tenantId, name: value})
}

func (c collection) CountDocuments(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return c.Collection.CountDocuments(ctx, filter)
}

func (c collection) FetchStream(ctx context.Context) (<-chan string, error) {
	ch := make(chan string)
	go c.consumeExistingEvents(ctx, ch)
	go c.consumeNewEvents(ctx, ch)
	go c.consumeErrorEvents(ctx, ch)
	return ch, nil
}

type Entry struct {
	Id string `bson:"id"`
}

func (c collection) consumeExistingEvents(ctx context.Context, ch chan<- string) {
	cursor, err := c.Collection.Find(ctx, bson.M{"status": bson.M{"$ne": "PROCESSED"}})
	if err != nil {
		log.Printf("Failed to find existing events: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry Entry
		if err := cursor.Decode(&entry); err != nil {
			log.Printf("Failed to decode existing entry: %v", err)
			continue
		}
		ch <- entry.Id
	}
}

func (c collection) consumeNewEvents(ctx context.Context, ch chan<- string) {
	pipeline := mongo.Pipeline{bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key:   "operationType",
			Value: "insert",
		}},
	}}}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := c.Collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Printf("Failed to start change stream: %v", err)
		return
	}
	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var changeEvent struct {
			FullDocument bson.M `bson:"fullDocument,omitempty"`
		}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change stream document: %v", err)
			continue
		}
		if id, ok := changeEvent.FullDocument["id"].(string); ok {
			ch <- id
		} else {
			log.Printf("Invalid _id type in change event")
		}
	}
}

func (c collection) consumeErrorEvents(ctx context.Context, ch chan<- string) {
	pipeline := mongo.Pipeline{bson.D{{
		Key: "$match", Value: bson.D{
			{
				Key:   "operationType",
				Value: "update",
			},
			{
				Key:   "fullDocument.status",
				Value: "ERROR",
			},
		},
	}}}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := c.Collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Printf("Failed to start change stream: %v", err)
		return
	}
	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var changeEvent struct {
			DocumentKey bson.M `bson:"documentKey,omitempty"`
		}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change stream document: %v", err)
			continue
		}
		if id, ok := changeEvent.DocumentKey["id"].(string); ok {
			go func(id string) {
				time.Sleep(5 * time.Second)
				ch <- id
			}(id)
		} else {
			log.Printf("Invalid _id type in change event")
		}
	}
}
