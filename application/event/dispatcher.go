package event

import (
	"bom-pedido-api/domain/events"
	"context"
	"os"
	"strconv"
)

type HandlerFunc func(message MessageEvent) error

type MessageEvent interface {
	Ack() error
	AckIfNoError(err error) error
	Nack()
	ParseData(event interface{})
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

type Emitter interface {
	Emit(ctx context.Context, event *events.Event) error
}

type Handler interface {
	Emitter
	Consume(options *ConsumerOptions, handler HandlerFunc)
	Close()
}
