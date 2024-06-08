package event

type EventHandler interface {
	Consume(topic string, handler func(event Event) error)
}
