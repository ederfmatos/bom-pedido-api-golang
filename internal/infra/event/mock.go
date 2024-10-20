package event

import (
	"bom-pedido-api/internal/application/event"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockEventHandler struct {
	mock.Mock
}

func NewMockEventHandler() *MockEventHandler {
	return &MockEventHandler{}
}

func (handler *MockEventHandler) Emit(ctx context.Context, event *event.Event) error {
	args := handler.Called(ctx, event)
	return args.Error(0)
}

func (handler *MockEventHandler) Close() {
}

func (handler *MockEventHandler) Consume(*event.ConsumerOptions, event.HandlerFunc) {
}
