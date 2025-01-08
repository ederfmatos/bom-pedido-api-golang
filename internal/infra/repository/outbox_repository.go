package repository

import (
	"bom-pedido-api/pkg/mongo"
	"context"
	"time"
)

type (
	Entry struct {
		Id              string     `json:"id,omitempty" bson:"_id"`
		Name            string     `json:"name,omitempty" bson:"name"`
		Data            string     `json:"data,omitempty" bson:"data"`
		CreatedAt       time.Time  `json:"createdAt,omitempty" bson:"createdAt"`
		Status          string     `json:"status,omitempty" bson:"status"`
		ProcessedAt     *time.Time `json:"processedAt" bson:"processedAt"`
		LastAttemptTime *time.Time `json:"lastAttemptTime" bson:"lastAttemptTime"`
	}

	MongoOutboxRepository struct {
		collection *mongo.Collection
	}
)

func NewMongoOutboxRepository(collection *mongo.Collection) *MongoOutboxRepository {
	return &MongoOutboxRepository{collection: collection}
}

func (r *MongoOutboxRepository) Save(ctx context.Context, entry *Entry) error {
	return r.collection.InsertOne(ctx, entry)
}

func (r *MongoOutboxRepository) Update(ctx context.Context, entry *Entry) error {
	return r.collection.UpdateByID(ctx, entry.Id, entry)
}

func (r *MongoOutboxRepository) Get(ctx context.Context, id string) (*Entry, error) {
	outbox := &Entry{}
	if err := r.collection.FindByID(ctx, id, outbox); err != nil {
		return nil, err
	}
	if outbox.Id == "" {
		return nil, nil
	}
	return outbox, nil
}

func (entry *Entry) MarkAsError() {
	now := time.Now()
	entry.LastAttemptTime = &now
	entry.Status = "ERROR"
}

func (entry *Entry) MarkAsProcessed() {
	now := time.Now()
	entry.ProcessedAt = &now
	entry.Status = "PROCESSED"
}
