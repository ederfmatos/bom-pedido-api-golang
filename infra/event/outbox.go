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

type OutboxEventHandler struct {
	handler          event.Handler
	outboxRepository outbox.Repository
	stream           Stream
}

func NewOutboxEventHandler(handler event.Handler, outboxRepository outbox.Repository, stream Stream) *OutboxEventHandler {
	eventHandler := &OutboxEventHandler{
		handler:          handler,
		outboxRepository: outboxRepository,
		stream:           stream,
	}
	eventHandler.handleStream()
	return eventHandler
}

func (handler *OutboxEventHandler) Emit(ctx context.Context, event *events.Event) error {
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
	return handler.outboxRepository.Save(ctx, entry)
}

func (handler *OutboxEventHandler) handleStream() {
	fetchEvents, err := handler.stream.FetchStream()
	if err != nil {
		panic(err)
	}
	go func() {
		for id := range fetchEvents {
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

func (handler *OutboxEventHandler) processEntry(entry *outbox.Entry) error {
	if entry == nil || entry.Status == "PROCESSED" {
		return nil
	}
	var messageEvent events.Event
	err := json.Unmarshal([]byte(entry.Data), &messageEvent)
	if err != nil {
		entry.MarkAsError()
		_ = handler.outboxRepository.Update(context.Background(), entry)
		return err
	}
	err = handler.handler.Emit(context.Background(), &messageEvent)
	if err != nil {
		entry.MarkAsError()
		_ = handler.outboxRepository.Update(context.Background(), entry)
		return err
	}
	entry.MarkAsProcessed()
	_ = handler.outboxRepository.Update(context.Background(), entry)
	return nil
}

func (handler *OutboxEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	handler.handler.Consume(options, handlerFunc)
}

func (handler *OutboxEventHandler) Close() {
	handler.handler.Close()
}
