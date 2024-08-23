package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/admin"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AdminSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	adminSqlRepository := NewDefaultAdminRepository(sqlConnection)
	runAdminTests(t, adminSqlRepository)
}

func Test_AdminMemoryRepository(t *testing.T) {
	adminSqlRepository := NewAdminMemoryRepository()
	runAdminTests(t, adminSqlRepository)
}

func runAdminTests(t *testing.T, repository repository.AdminRepository) {
	ctx := context.TODO()

	aAdmin, err := admin.New(faker.Name(), faker.Email(), faker.Word())
	assert.NoError(t, err)

	savedAdmin, err := repository.FindByEmail(ctx, aAdmin.GetEmail())
	assert.NoError(t, err)
	assert.Nil(t, savedAdmin)

	err = repository.Create(ctx, aAdmin)
	assert.NoError(t, err)

	savedAdmin, err = repository.FindByEmail(ctx, aAdmin.GetEmail())
	assert.NoError(t, err)
	assert.NotNil(t, savedAdmin)
	assert.Equal(t, aAdmin, savedAdmin)
}
