package event

import (
	"context"
	"os"
	"strconv"
)

type Handler func(message MessageEvent) error

type MessageEvent interface {
	Ack() error
	Nack()
	AsEvent() *Event
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
	WorkerPoolSize int
}

func NewConsumerOptions(id string, workerPoolSize int) *ConsumerOptions {
	return &ConsumerOptions{Id: id, WorkerPoolSize: workerPoolSize}
}

func OptionsForQueue(queue string) *ConsumerOptions {
	return NewConsumerOptions(queue, defaultWorkerPoolSize)
}

type Dispatcher interface {
	Emit(ctx context.Context, event *Event) error
	Consume(options *ConsumerOptions, handler Handler)
	Close()
}
