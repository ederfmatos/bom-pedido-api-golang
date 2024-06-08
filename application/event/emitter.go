package event

type EventEmitter interface {
	Emit(event *Event) error
}
