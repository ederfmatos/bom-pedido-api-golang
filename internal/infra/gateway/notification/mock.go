package notification

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity/notification"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockNotificationGateway struct {
	mock.Mock
}

func (f *MockNotificationGateway) Name() string {
	return "MOCK"
}

func NewMockNotificationGateway() gateway.NotificationGateway {
	return &MockNotificationGateway{}
}

func (f *MockNotificationGateway) Send(ctx context.Context, notification *notification.Notification) error {
	args := f.Called(ctx, notification)
	return args.Error(0)
}
