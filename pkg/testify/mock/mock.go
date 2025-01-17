package mock

import (
	"github.com/stretchr/testify/mock"
	"runtime"
	"strings"
)

const Anything = mock.Anything

type Mock struct {
	mock.Mock
}

func (m *Mock) Called(arguments ...interface{}) mock.Arguments {
	pc, _, _, _ := runtime.Caller(1)
	functionPath := runtime.FuncForPC(pc).Name()
	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]
	return m.MethodCalled(functionName, arguments...)
}
