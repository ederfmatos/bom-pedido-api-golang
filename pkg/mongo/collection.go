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

type Collection struct {
	*mongo.Collection
}

type Database struct {
	database *mongo.Database
}

func NewDatabase(database *mongo.Database) *Database {
	return &Database{database: database}
}

func (d Database) ForCollection(name string) *Collection {
	return newCollection(d.database.Collection(name))
}

func newCollection(collection *mongo.Collection) *Collection {
	return &Collection{Collection: collection}
}

// Private reusable method for find operations
func (c Collection) findOne(ctx context.Context, filter bson.M, target any) error {
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

func (c Collection) DeleteByID(ctx context.Context, id string) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

func (c Collection) Upsert(ctx context.Context, id string, value any) error {
	return c.update(ctx, id, value, true)
}

func (c Collection) UpdateByID(ctx context.Context, id string, value any) error {
	return c.update(ctx, id, value, false)
}

func (c Collection) update(ctx context.Context, id string, value any, upsert bool) error {
	update := bson.M{"$set": value}
	updateOptions := options.Update().SetUpsert(upsert)
	filter := bson.D{{Key: "id", Value: id}}
	_, err := c.Collection.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return fmt.Errorf("update record: %w", err)
	}
	return nil
}

func (c Collection) FindByID(ctx context.Context, id string, target any) error {
	return c.findOne(ctx, bson.M{"id": id}, target)
}

func (c Collection) FindByTenantIdAnd(ctx context.Context, tenantId, param, value string, target any) error {
	return c.findOne(ctx, bson.M{"tenantId": tenantId, param: value}, target)
}

func (c Collection) FindBy(ctx context.Context, param, value string, target any) error {
	return c.findOne(ctx, bson.M{param: value}, target)
}

func (c Collection) FindByValues(ctx context.Context, values map[string]interface{}, target any) error {
	return c.findOne(ctx, values, target)
}

func (c Collection) FindAllByID(ctx context.Context, ids []string) (*mongo.Cursor, error) {
	filter := bson.M{"id": bson.M{"$in": ids}}
	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find records: %w", err)
	}
	return cursor, nil
}

func (c Collection) FindAllBy(ctx context.Context, values map[string]interface{}) (*mongo.Cursor, error) {
	filter := bson.M(values)
	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find records: %w", err)
	}
	return cursor, nil
}

func (c Collection) InsertOne(ctx context.Context, value any) error {
	_, err := c.Collection.InsertOne(ctx, value)
	if err != nil {
		return fmt.Errorf("insert record: %w", err)
	}
	return nil
}

func (c Collection) ExistsByID(ctx context.Context, id string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{"id": id})
}

func (c Collection) ExistsBy(ctx context.Context, name string, value string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{name: value})
}

func (c Collection) ExistsByTenantIdAnd(ctx context.Context, tenantId, name, value string) (bool, error) {
	return c.existsByValues(ctx, map[string]interface{}{"tenantId": tenantId, name: value})
}

func (c Collection) existsByValues(ctx context.Context, values map[string]interface{}) (bool, error) {
	result := c.Collection.FindOne(ctx, bson.M(values))
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("find record: %w", err)
	}
	return true, nil
}

func (c Collection) FetchStream(ctx context.Context) (<-chan string, error) {
	ch := make(chan string)
	go c.consumeExistingEvents(ctx, ch)
	go c.consumeNewEvents(ctx, ch)
	go c.consumeErrorEvents(ctx, ch)
	return ch, nil
}

type Entry struct {
	Id string `bson:"_id"`
}

func (c Collection) consumeExistingEvents(ctx context.Context, ch chan<- string) {
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

func (c Collection) consumeNewEvents(ctx context.Context, ch chan<- string) {
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
			DocumentKey bson.M `bson:"documentKey,omitempty"`
		}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change stream document: %v", err)
			continue
		}
		if id, ok := changeEvent.DocumentKey["_id"].(string); ok {
			ch <- id
		} else {
			log.Printf("Invalid _id type in change event")
		}
	}
}

func (c Collection) consumeErrorEvents(ctx context.Context, ch chan<- string) {
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
		if id, ok := changeEvent.DocumentKey["_id"].(string); ok {
			go func(id string) {
				time.Sleep(5 * time.Second)
				ch <- id
			}(id)
		} else {
			log.Printf("Invalid _id type in change event")
		}
	}
}
