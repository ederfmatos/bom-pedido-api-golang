package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/transaction"
	"context"
)

const (
	sqlCreatePixTransaction       = "INSERT INTO pix_transactions (id, order_id, amount, status, qr_code, payment_gateway, qr_code_link, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	sqlExistsTransactionByOrderId = "SELECT 1 FROM transactions WHERE order_id = $1 LIMIT 1"
	pix                           = "PIX"
)

type DefaultTransactionRepository struct {
	SqlConnection
}

func NewDefaultTransactionRepository(sqlConnection SqlConnection) repository.TransactionRepository {
	return &DefaultTransactionRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultTransactionRepository) CreatePixTransaction(ctx context.Context, transaction *transaction.PixTransaction) error {
	return repository.Sql(sqlCreatePixTransaction).
		Values(transaction.Id, transaction.OrderId, transaction.Amount, transaction.Status, transaction.QrCode, transaction.PaymentGateway, transaction.QrCodeLink, pix).
		Update(ctx)
}

func (repository *DefaultTransactionRepository) ExistsByOrderId(ctx context.Context, id string) (bool, error) {
	return repository.Sql(sqlExistsTransactionByOrderId).Values(id).Exists(ctx)
}
