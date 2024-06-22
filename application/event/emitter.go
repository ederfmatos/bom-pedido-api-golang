package event

import "context"

type EventEmitter interface {
	Emit(ctx context.Context, event *Event) error
}
