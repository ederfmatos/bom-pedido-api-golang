package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/internal/infra/retry"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
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
	ctx, span := telemetry.StartSpan(ctx, "OutboxEventEmitter.Emit",
		"event", event.Name,
		"eventId", event.Id,
		"eventCorrelationId", event.CorrelationId,
	)
	defer span.End()
	eventData, err := json.Marshal(ctx, event)
	if err != nil {
		return err
	}
	entry := &repository.Entry{
		Id:        event.Id,
		Name:      event.Name,
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	ctx, span := telemetry.StartSpan(ctx, "OutboxEventHandler.Process", "eventId", id)
	defer span.End()
	defer cancel()
	_ = handler.locker.LockFunc(ctx, id, time.Minute, func() {
		entry, err := handler.outboxRepository.Get(ctx, id)
		if err != nil {
			return
		}
		err = retry.Func(ctx, 5, time.Second, time.Second*30, func(ctx context.Context) error {
			return handler.processEntry(ctx, entry)
		})
		span.RecordError(err)
	})
}

func (handler *OutboxEventHandler) processEntry(ctx context.Context, entry *repository.Entry) error {
	if entry == nil || entry.Status == "PROCESSED" {
		return nil
	}
	var messageEvent event.Event
	err := json.Unmarshal(ctx, []byte(entry.Data), &messageEvent)
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

func (handler *OutboxEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	handler.handler.Consume(options, handlerFunc)
}

func (handler *OutboxEventHandler) Close() {
	handler.handler.Close()
}
