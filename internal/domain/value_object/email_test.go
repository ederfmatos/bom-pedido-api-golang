package value_object

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		email       string
		expectedErr string
	}{
		{"example@example.com", ""},
		{"invalid-email", "invalid email"},
		{"another.valid@example.com", ""},
		{"", "email is empty"},
		{"toolong" + string(make([]byte, 250)) + "@example.com", "email is too long"},
	}

	for _, item := range tests {
		_, err := NewEmail(item.email)
		if err != nil && err.Error() != item.expectedErr {
			t.Errorf("expected error '%v', got '%v'", item.expectedErr, err)
		} else if err == nil && item.expectedErr != "" {
			t.Errorf("expected error '%v', got nil", item.expectedErr)
		}
	}
}
