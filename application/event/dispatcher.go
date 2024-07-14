package event

import (
	"context"
	"os"
	"strconv"
)

type HandlerFunc func(message MessageEvent) error

type MessageEvent interface {
	Ack() error
	AckIfNoError(err error) error
	Nack()
	GetEvent() *Event
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
