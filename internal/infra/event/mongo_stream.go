package event

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoStream struct {
	collection *mongo.Collection
}

func NewMongoStream(collection *mongo.Collection) Stream {
	return &MongoStream{collection: collection}
}

func (stream *MongoStream) FetchStream() (chan string, error) {
	ch := make(chan string)
	go stream.consumeExistingEvents(ch)
	go stream.consumeErrorEvents(ch)
	go stream.consumeNewEvents(ch)
	return ch, nil
}

type Entry struct {
	Id string `bson:"_id"`
}

func (stream *MongoStream) consumeExistingEvents(ch chan string) {
	ctx := context.Background()
	cursor, err := stream.collection.Find(ctx, bson.M{"status": bson.M{"$ne": "PROCESSED"}})
	if err != nil {
		log.Fatalf("Failed to find existing events: %v", err)
	}

	for cursor.Next(ctx) {
		var entry Entry
		if err := cursor.Decode(&entry); err != nil {
			log.Printf("Failed to decode existing entry: %v", err)
			continue
		}
		ch <- entry.Id
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
	_ = cursor.Close(ctx)
}

func (stream *MongoStream) consumeNewEvents(ch chan string) {
	ctx := context.Background()
	pipeline := mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "operationType", Value: "insert"}}}}}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := stream.collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Fatalf("Failed to start change stream: %v", err)
	}

	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var changeEvent struct {
			DocumentKey primitive.M `bson:"documentKey,omitempty"`
		}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change stream document: %v", err)
			continue
		}
		ch <- changeEvent.DocumentKey["_id"].(string)
	}
}

func (stream *MongoStream) consumeErrorEvents(ch chan string) {
	ctx := context.Background()
	pipeline := mongo.Pipeline{bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "update"},
			{Key: "fullDocument.status", Value: "ERROR"},
		}},
	}}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := stream.collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Fatalf("Failed to start change stream: %v", err)
	}

	for changeStream.Next(ctx) {
		var changeEvent struct {
			DocumentKey primitive.M `bson:"documentKey,omitempty"`
		}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change stream document: %v", err)
			continue
		}
		go func() {
			time.Sleep(time.Second * 5)
			ch <- changeEvent.DocumentKey["_id"].(string)
		}()
	}
	_ = changeStream.Close(ctx)
}
