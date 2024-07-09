package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/domain/events"
	"bom-pedido-api/infra/repository/outbox"
	"bom-pedido-api/infra/retry"
	"context"
	"encoding/json"
	"time"
)

type DynamoStreamsEventHandler struct {
	handler          event.Handler
	outboxRepository outbox.Repository
	dynamoStream     *DynamoStream
}

func NewDynamoStreamsEventHandler(handler event.Handler, outboxRepository outbox.Repository, dynamoStream *DynamoStream) *DynamoStreamsEventHandler {
	dynamoEventHandler := &DynamoStreamsEventHandler{
		handler:          handler,
		outboxRepository: outboxRepository,
		dynamoStream:     dynamoStream,
	}
	dynamoEventHandler.handleStream()
	return dynamoEventHandler
}

func (handler *DynamoStreamsEventHandler) Emit(ctx context.Context, event *events.Event) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	entry := &outbox.Entry{
		Id:        event.Id,
		Name:      event.Name,
		Data:      string(eventData),
		CreatedAt: time.Now(),
		Status:    "NEW",
	}
	return handler.outboxRepository.Store(ctx, entry)
}

func (handler *DynamoStreamsEventHandler) handleStream() {
	fetchEvents, err := handler.dynamoStream.FetchEvents()
	if err != nil {
		panic(err)
	}
	go func() {
		for record := range fetchEvents {
			id := *record.Dynamodb.NewImage["id"].S
			entry, err := handler.outboxRepository.Get(context.Background(), id)
			if err != nil {
				continue
			}
			retryable := retry.NewRetry(5, time.Second*2, time.Minute)
			go retryable.Execute(func() error {
				return handler.processEntry(entry)
			})
		}
	}()
}

func (handler *DynamoStreamsEventHandler) processEntry(entry *outbox.Entry) error {
	if entry == nil || entry.Status == "PROCESSED" {
		return nil
	}
	var messageEvent events.Event
	err := json.Unmarshal([]byte(entry.Data), &messageEvent)
	if err != nil {
		_ = handler.outboxRepository.MarkAsError(context.Background(), entry)
		return err
	}
	err = handler.handler.Emit(context.Background(), &messageEvent)
	if err != nil {
		_ = handler.outboxRepository.MarkAsError(context.Background(), entry)
		return err
	}
	_ = handler.outboxRepository.MarkAsProcessed(context.Background(), entry)
	return nil
}

func (handler *DynamoStreamsEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	handler.handler.Consume(options, handlerFunc)
}

func (handler *DynamoStreamsEventHandler) Close() {
	handler.handler.Close()
}
