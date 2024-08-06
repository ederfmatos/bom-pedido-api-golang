package repository

import (
	"bom-pedido-api/infra/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	container := test.NewContainer()
	defer container.Down()
	code := m.Run()
	os.Exit(code)
}
