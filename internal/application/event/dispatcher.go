package event

import (
	"context"
)

type HandlerFunc func(ctx context.Context, message *MessageEvent) error

type MessageEvent struct {
	Event  *Event
	AckFn  func(context.Context) error
	NackFn func(context.Context)
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

func (m *MessageEvent) Nack(ctx context.Context) {
	m.NackFn(ctx)
}

type Emitter interface {
	Emit(ctx context.Context, event *Event) error
}

type Handler interface {
	Emitter
	OnEvent(eventName string, handler HandlerFunc)
	Close()
}
