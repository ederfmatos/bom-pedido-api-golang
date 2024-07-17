package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/value_object"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OrderSqlRepository(t *testing.T) {
	database, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Error(err)
	}
	defer database.Close()
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS orders
		(
			id                VARCHAR(36) NOT NULL PRIMARY KEY,
			code              INTEGER NULL DEFAULT 1,
			customer_id       VARCHAR(36) NOT NULL,
			payment_method    VARCHAR(30) NOT NULL,
			payment_mode      VARCHAR(30) NOT NULL,
			delivery_mode     VARCHAR(30) NOT NULL,
			status            VARCHAR(30) NOT NULL,
			credit_card_token VARCHAR(255),
			'change'          DECIMAL(6, 2),
			delivery_time     TIMESTAMP   NOT NULL,
			created_at        TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS order_items
		(
			id          VARCHAR(36) NOT NULL PRIMARY KEY,
			order_id    VARCHAR(36) NOT NULL,
			product_id  VARCHAR(36) NOT NULL,
			status      VARCHAR(30) NOT NULL,
			quantity    NUMERIC     NOT NULL,
			observation TEXT,
			price       DECIMAL(6, 2),
			created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS order_history
		(
			id         SERIAL      PRIMARY KEY,
			order_id   VARCHAR(36) NOT NULL,
			changed_by VARCHAR(36) NOT NULL,
			changed_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status     VARCHAR(30) NOT NULL,
			data       TEXT
		);
	`)
	assert.NoError(t, err)
	sqlConnection := NewDefaultSqlConnection(database)
	orderRepository := NewDefaultOrderRepository(sqlConnection)
	orderTests(t, orderRepository)
}

func Test_OrderMemoryRepository(t *testing.T) {
	orderRepository := NewOrderMemoryRepository()
	orderTests(t, orderRepository)
}

func orderTests(t *testing.T, repository repository.OrderRepository) {
	ctx := context.TODO()

	customerId := value_object.NewID()
	adminId := value_object.NewID()
	anOrder, err := order.New(customerId, enums.CreditCard, enums.InReceiving, enums.Delivery, "", 10.0, time.Now())
	anOrder.Code = 1
	assert.NoError(t, err)

	savedOrder, err := repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Nil(t, savedOrder)

	err = repository.Create(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)

	// Approve
	err = anOrder.Approve(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)

	// MarkAsInProgress
	err = anOrder.MarkAsInProgress(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)

	// MarkAsAwaitingDelivery
	err = anOrder.MarkAsAwaitingDelivery(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)

	// MarkAsDelivering
	err = anOrder.MarkAsDelivering(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)

	// Finish
	err = anOrder.Finish(time.Now(), adminId)
	assert.NoError(t, err)

	err = repository.Update(ctx, anOrder)
	assert.NoError(t, err)

	savedOrder, err = repository.FindById(ctx, anOrder.Id)
	assert.NoError(t, err)
	assert.Equal(t, *anOrder, *savedOrder)
}
