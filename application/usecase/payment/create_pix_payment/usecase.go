package create_pix_payment

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/application/repository"
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
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() {
		return err
	}
	aCustomer, err := uc.customerRepository.FindById(ctx, anOrder.CustomerID)
	if err != nil || aCustomer == nil {
		return err
	}
	existsTransaction, err := uc.transactionRepository.ExistsByOrderId(ctx, anOrder.Id)
	if err != nil || existsTransaction {
		return err
	}
	createPixInput := gateway.CreateQrCodePixInput{
		InternalOrderId: anOrder.Id,
		Amount:          anOrder.Amount,
		MerchantId:      anOrder.MerchantId,
		Description:     "Pedido no Bom Pedido",
		Customer: gateway.PixCustomer{
			Name:  aCustomer.Name,
			Email: aCustomer.GetEmail(),
		},
	}
	createPixOutput, err := uc.pixGateway.CreateQrCodePix(ctx, createPixInput)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentCreated(anOrder.Id, createPixOutput.Id, createPixOutput.PaymentGateway))
}
