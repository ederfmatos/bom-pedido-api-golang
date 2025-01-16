package require

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}

func Equal(t *testing.T, expected, actual interface{}) {
	require.Equal(t, expected, actual)
}

func ErrorIs(t *testing.T, err error, target error, msgAndArgs ...interface{}) {
	require.ErrorIs(t, err, target, msgAndArgs...)
}

func Error(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.Error(t, err, msgAndArgs...)
}

func NotNil(t *testing.T, object interface{}) {
	require.NotNil(t, object)
}

func Nil(t *testing.T, object interface{}) {
	require.Nil(t, object)
}

func False(t *testing.T, value bool) {
	require.False(t, value)
}

func True(t *testing.T, value bool) {
	require.True(t, value)
}

func Empty(t *testing.T, value interface{}) {
	require.Empty(t, value)
}

func NotEmpty(t *testing.T, value interface{}) {
	require.NotEmpty(t, value)
}

func ErrorAs(t *testing.T, err error, target interface{}, msgAndArgs ...interface{}) {
	require.ErrorAs(t, err, target, msgAndArgs...)
}
