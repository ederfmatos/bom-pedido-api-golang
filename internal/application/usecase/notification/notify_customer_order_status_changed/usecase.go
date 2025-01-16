package notify_customer_order_status_changed

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/entity/status"
	"context"
)

type (
	UseCase struct {
		orderRepository                repository.OrderRepository
		customerNotificationRepository repository.CustomerNotificationRepository
		notificationRepository         repository.NotificationRepository
	}

	Input struct {
		OrderId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository:                factory.OrderRepository,
		customerNotificationRepository: factory.CustomerNotificationRepository,
		notificationRepository:         factory.NotificationRepository,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) error {
	order, err := u.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil {
		return err
	}
	title, body := u.getNotification(order)
	if title == "" || body == "" {
		return nil
	}
	customerNotification, err := u.customerNotificationRepository.FindByCustomerId(ctx, order.CustomerID)
	if err != nil || customerNotification == nil {
		return err
	}
	notification := entity.NewNotification(title, body, customerNotification.Recipient, order.Id)
	notification.Put("orderId", order.Id).
		Put("event", "ORDER_STATUS_CHANGED").
		Put("status", order.GetStatus())
	return u.notificationRepository.Create(ctx, notification)
}

func (u *UseCase) getNotification(order *entity.Order) (string, string) {
	switch order.GetStatus() {

	case status.AwaitingApprovalStatus.Name():
		return "Pedido aguardando aprovação", "Seu pedido está aguardando aprovação, assim que tivermos uma atualização te notificamos"

	case status.ApprovedStatus.Name():
		return "Pedido aprovado", "Seu pedido foi aprovado"

	case status.InProgressStatus.Name():
		return "Pedido em preparação", "Seu pedido está sendo preparado"

	case status.RejectedStatus.Name():
		return "Pedido Rejeitado", "Seu pedido foi rejeitado. Clique aqui para ver mais detalhes."

	case status.CancelledStatus.Name():
		return "Pedido cancelado", "Seu pedido foi cancelado. Clique aqui para ver mais detalhes."

	case status.DeliveringStatus.Name():
		return "Pedido em rota de entrega", "Seu pedido está em rota de entrega"

	case status.AwaitingWithdrawStatus.Name():
		return "Pedido pronto para retirada", "Seu pedido está pronto, pode retira-lo no estabelecimento."

	case status.AwaitingDeliveryStatus.Name():
		return "Pedido aguardando entrega", "Seu pedido está pronto, só aguardando o entregador sair para entrega"

	case status.FinishedStatus.Name():
		return "Pedido finalizado", "Seu pedido foi finalizado. Acesso os detalhes do pedido para avaliar"

	default:
		return "", ""
	}
}
