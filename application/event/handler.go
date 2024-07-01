package event

type EventHandler interface {
	Consume(queue string, handler func(event Event) error)
	Close()
}
