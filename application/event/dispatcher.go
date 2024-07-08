package event

import "context"

type Handler func(event Event) error

type Dispatcher interface {
	Emit(ctx context.Context, event *Event) error
	Consume(id string, handler Handler)
	Close()
}
