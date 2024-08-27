package event

import (
	"bom-pedido-api/infra/telemetry"
	"context"
	"os"
	"strconv"
)

type HandlerFunc func(ctx context.Context, message *MessageEvent) error

type MessageEvent struct {
	Event  *Event
	AckFn  func(context.Context) error
	NackFn func(context.Context)
}

func (m *MessageEvent) Ack(ctx context.Context) error {
	_, span := telemetry.StartSpan(ctx, "MessageEvent.Ack")
	defer span.End()
	return m.AckFn(ctx)
}

func (m *MessageEvent) AckIfNoError(ctx context.Context, err error) error {
	if err == nil {
		return m.Ack(ctx)
	}
	return err
}

func (m *MessageEvent) Nack(ctx context.Context) {
	_, span := telemetry.StartSpan(ctx, "MessageEvent.Nack")
	defer span.End()
	m.NackFn(ctx)
}

var defaultWorkerPoolSize int

func init() {
	size := os.Getenv("POOL_WORKER_SIZE")
	if value, err := strconv.Atoi(size); err != nil {
		defaultWorkerPoolSize = 20
	} else {
		defaultWorkerPoolSize = value
	}
}

type ConsumerOptions struct {
	Id             string
	Queue          string
	Topics         []string
	WorkerPoolSize int
}

func OptionsForTopics(id string, topics ...string) *ConsumerOptions {
	return &ConsumerOptions{
		Id:             id,
		Queue:          id,
		Topics:         topics,
		WorkerPoolSize: defaultWorkerPoolSize,
	}
}

type Emitter interface {
	Emit(ctx context.Context, event *Event) error
}

type Handler interface {
	Emitter
	Consume(options *ConsumerOptions, handler HandlerFunc)
	Close()
}
