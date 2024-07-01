package repository

import (
	"bom-pedido-api/domain/entity"
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultProductRepository_Create(t *testing.T) {
	database, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Error(err)
	}
	defer database.Close()
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS products
		(
			id          VARCHAR(36)                            NOT NULL PRIMARY KEY,
			name        VARCHAR(255)                           NOT NULL,
			description MEDIUMTEXT,
			price       DECIMAL(6, 2)                          NOT NULL,
			status      VARCHAR(20) NOT NULL,
			created_at  TIMESTAMP                              NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	assert.NoError(t, err)
	sqlConnection := NewDefaultSqlConnection(database)
	productRepository := NewDefaultProductRepository(sqlConnection)
	ctx := context.Background()

	t.Run("should create a product", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0)
		assert.NoError(t, err)

		err = productRepository.Create(ctx, product)
		assert.NoError(t, err)

		savedProduct, err := productRepository.FindById(ctx, product.Id)
		assert.NoError(t, err)
		assert.Equal(t, product.Id, savedProduct.Id)
		assert.Equal(t, product.Name, savedProduct.Name)
		assert.Equal(t, product.Description, savedProduct.Description)
		assert.Equal(t, product.Price, savedProduct.Price)
		assert.Equal(t, product.Status, savedProduct.Status)
	})

	t.Run("should not create duplicated product", func(t *testing.T) {
		product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0)
		assert.NoError(t, err)

		err = productRepository.Create(ctx, product)
		assert.NoError(t, err)

		err = productRepository.Create(ctx, product)
		assert.Error(t, err)
	})
}
