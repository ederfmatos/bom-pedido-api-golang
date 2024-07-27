package event

import (
	"context"
	"os"
	"strconv"
)

type HandlerFunc func(message *MessageEvent) error

type MessageEvent struct {
	AckFn      func() error
	NackFn     func()
	GetEventFn func() *Event
}

func (m *MessageEvent) Ack() error {
	return m.AckFn()
}

func (m *MessageEvent) AckIfNoError(err error) error {
	if err == nil {
		return m.Ack()
	}
	return err
}

func (m *MessageEvent) NackIfError(err error) {
	if err != nil {
		m.Nack()
	}
}

func (m *MessageEvent) Nack() {
	m.NackFn()
}

func (m *MessageEvent) GetEvent() *Event {
	return m.GetEventFn()
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
