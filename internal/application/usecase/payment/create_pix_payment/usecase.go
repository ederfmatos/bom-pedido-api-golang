package create_pix_payment

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		customerRepository    repository.CustomerRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	Input struct {
		OrderId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository:       factory.OrderRepository,
		customerRepository:    factory.CustomerRepository,
		transactionRepository: factory.TransactionRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "CREATE_PIX_PAYMENT_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	order, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsPixInApp() {
		return err
	}
	customer, err := uc.customerRepository.FindById(ctx, order.CustomerID)
	if err != nil || customer == nil {
		return err
	}
	existsTransaction, err := uc.transactionRepository.ExistsByOrderId(ctx, order.Id)
	if err != nil || existsTransaction {
		return err
	}
	createPixInput := gateway.CreateQrCodePixInput{
		InternalOrderId: order.Id,
		Amount:          order.Amount,
		MerchantId:      order.MerchantId,
		Description:     "Pedido no Bom Pedido",
		Customer: gateway.PixCustomer{
			Name:  customer.Name,
			Email: customer.GetEmail(),
		},
	}
	createPixOutput, err := uc.pixGateway.CreateQrCodePix(ctx, createPixInput)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentCreated(order.Id, createPixOutput.Id, createPixOutput.PaymentGateway))
}
