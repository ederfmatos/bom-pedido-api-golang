package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/internal/infra/retry"
	"context"
	"encoding/json"
	"time"
)

type outboxRepository interface {
	Save(ctx context.Context, entry *repository.Entry) error
	Get(ctx context.Context, id string) (*repository.Entry, error)
	Update(ctx context.Context, entry *repository.Entry) error
	FetchStream(ctx context.Context) (<-chan string, error)
}

type OutboxEventHandler struct {
	handler          event.Handler
	outboxRepository outboxRepository
	locker           lock.Locker
}

func NewOutboxEventHandler(handler event.Handler, outboxRepository outboxRepository, locker lock.Locker) (*OutboxEventHandler, error) {
	eventHandler := &OutboxEventHandler{
		handler:          handler,
		outboxRepository: outboxRepository,
		locker:           locker,
	}
	if err := eventHandler.handleStream(); err != nil {
		return nil, err
	}
	return eventHandler, nil
}

func (handler *OutboxEventHandler) Emit(ctx context.Context, event *event.Event) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	entry := &repository.Entry{
		Id:        event.Id,
		Name:      string(event.Name),
		Data:      string(eventData),
		CreatedAt: time.Now(),
		Status:    "NEW",
	}
	return handler.outboxRepository.Save(ctx, entry)
}

func (handler *OutboxEventHandler) handleStream() error {
	fetchEvents, err := handler.outboxRepository.FetchStream(context.Background())
	if err != nil {
		return err
	}
	go func() {
		for id := range fetchEvents {
			handler.processEvent(id)
		}
	}()
	return nil
}

func (handler *OutboxEventHandler) processEvent(id string) {
	// TODO: Add telemetry
	ctx := context.Background()
	_ = handler.locker.LockFunc(ctx, id, func() {
		entry, err := handler.outboxRepository.Get(ctx, id)
		if err != nil || entry == nil {
			return
		}
		_ = retry.Func(ctx, 5, time.Second, time.Second*30, func(ctx context.Context) error {
			return handler.processEntry(ctx, entry)
		})
	})
}

func (handler *OutboxEventHandler) processEntry(ctx context.Context, entry *repository.Entry) error {
	if entry == nil || entry.Status == "PROCESSED" {
		return nil
	}
	var messageEvent event.Event
	err := json.Unmarshal([]byte(entry.Data), &messageEvent)
	if err != nil {
		entry.MarkAsError()
		_ = handler.outboxRepository.Update(ctx, entry)
		return err
	}
	err = handler.handler.Emit(ctx, &messageEvent)
	if err != nil {
		entry.MarkAsError()
		_ = handler.outboxRepository.Update(ctx, entry)
		return err
	}
	entry.MarkAsProcessed()
	_ = handler.outboxRepository.Update(ctx, entry)
	return nil
}

func (handler *OutboxEventHandler) OnEvent(event string, handlerFunc event.HandlerFunc) {
	handler.handler.OnEvent(event, handlerFunc)
}

func (handler *OutboxEventHandler) Close() {
	handler.handler.Close()
}
