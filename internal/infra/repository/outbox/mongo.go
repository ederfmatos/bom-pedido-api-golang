package outbox

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOutboxRepository struct {
	collection *mongo.Collection
}

func NewMongoOutboxRepository(collection *mongo.Collection) *MongoOutboxRepository {
	return &MongoOutboxRepository{collection: collection}
}

func (r *MongoOutboxRepository) Save(ctx context.Context, entry *Entry) error {
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

func (r *MongoOutboxRepository) Update(ctx context.Context, entry *Entry) error {
	update := bson.M{"$set": entry}
	_, err := r.collection.UpdateByID(ctx, entry.Id, update)
	return err
}

func (r *MongoOutboxRepository) Get(ctx context.Context, id string) (*Entry, error) {
	result := r.collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}
	outbox := &Entry{}
	err := result.Decode(outbox)
	if err != nil {
		return nil, err
	}
	return outbox, nil
}
