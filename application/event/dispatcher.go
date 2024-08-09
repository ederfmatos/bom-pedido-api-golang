package event

import (
	"context"
	"os"
	"strconv"
)

type HandlerFunc func(ctx context.Context, message *MessageEvent) error

type MessageEvent struct {
	AckFn      func(context.Context) error
	NackFn     func(context.Context)
	GetEventFn func(context.Context) *Event
}

func (m *MessageEvent) Ack(ctx context.Context) error {
	return m.AckFn(ctx)
}

func (m *MessageEvent) AckIfNoError(ctx context.Context, err error) error {
	if err == nil {
		return m.Ack(ctx)
	}
	return err
}

func (m *MessageEvent) NackIfError(ctx context.Context, err error) {
	if err != nil {
		m.Nack(ctx)
	}
}

func (m *MessageEvent) Nack(ctx context.Context) {
	m.NackFn(ctx)
}

func (m *MessageEvent) GetEvent(ctx context.Context) *Event {
	return m.GetEventFn(ctx)
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
	TopicName      string
	WorkerPoolSize int
}

func NewConsumerOptions(queue, eventName, id string, workerPoolSize int) *ConsumerOptions {
	return &ConsumerOptions{Queue: queue, TopicName: eventName, Id: id, WorkerPoolSize: workerPoolSize}
}

func OptionsForTopic(topicName, id string) *ConsumerOptions {
	return &ConsumerOptions{
		Id:             id,
		TopicName:      topicName,
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
