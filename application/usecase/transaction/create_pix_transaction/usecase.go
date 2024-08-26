package create_pix_transaction

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/transaction"
	"bom-pedido-api/infra/retry"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		customerRepository    repository.CustomerRepository
		transactionRepository repository.TransactionRepository
		merchantRepository    repository.MerchantRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
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
		merchantRepository:    factory.MerchantRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
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
	createPixOutput, err := uc.createQrCodePix(ctx, anOrder, aCustomer)
	if err != nil {
		return err
	}
	err = uc.createPixTransaction(ctx, createPixOutput, anOrder)
	if err != nil {
		_ = uc.eventEmitter.Emit(ctx, event.NewRefundTransactionEvent(anOrder.Id, createPixOutput.Id))
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionCreated(anOrder.Id, createPixOutput.Id))
}

func (uc *UseCase) createPixTransaction(ctx context.Context, pix *gateway.CreateQrCodePixOutput, order *order.Order) error {
	pixTransaction := transaction.NewPixTransaction(pix.Id, order.Id, pix.QrCode, pix.PaymentGateway, pix.QrCodeLink, order.Amount)
	return retry.Func(ctx, 5, time.Second, time.Second*10, func(ctx context.Context) error {
		return uc.transactionRepository.CreatePixTransaction(ctx, pixTransaction)
	})
}

func (uc *UseCase) createQrCodePix(ctx context.Context, anOrder *order.Order, customer *customer.Customer) (*gateway.CreateQrCodePixOutput, error) {
	createPixInput := gateway.CreateQrCodePixInput{
		Amount:          anOrder.Amount,
		InternalOrderId: anOrder.Id,
		Description:     "Pedido no Bom Pedido",
		MerchantId:      anOrder.MerchantId,
		Customer: gateway.PixCustomer{
			Name:  customer.Name,
			Email: customer.GetEmail(),
		},
	}
	return uc.pixGateway.CreateQrCodePix(ctx, createPixInput)
}
