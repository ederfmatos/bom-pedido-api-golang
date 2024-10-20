package repository

import (
	"bom-pedido-api/internal/infra/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	container := test.NewContainer()
	defer container.Down()
	code := m.Run()
	os.Exit(code)
}
