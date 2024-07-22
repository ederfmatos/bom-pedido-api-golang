package outbox

import (
	"bom-pedido-api/infra/telemetry"
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
	ctx, span := telemetry.StartSpan(ctx, "MongoOutboxRepository.Save", "entry", entry.Id)
	defer span.End()
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

func (r *MongoOutboxRepository) Update(ctx context.Context, entry *Entry) error {
	ctx, span := telemetry.StartSpan(ctx, "MongoOutboxRepository.Update", "entry", entry.Id)
	defer span.End()
	update := bson.M{"$set": entry}
	_, err := r.collection.UpdateByID(ctx, entry.Id, update)
	return err
}

func (r *MongoOutboxRepository) Get(ctx context.Context, id string) (*Entry, error) {
	ctx, span := telemetry.StartSpan(ctx, "MongoOutboxRepository.Get", "entry", id)
	defer span.End()
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
