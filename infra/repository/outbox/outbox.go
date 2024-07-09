package outbox

import (
	"context"
	"time"
)

type Entry struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Data      string    `json:"data,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Status    string    `json:"status,omitempty"`
}

type Repository interface {
	Store(ctx context.Context, entry *Entry) error
	Get(ctx context.Context, id string) (*Entry, error)
	MarkAsProcessed(ctx context.Context, entry *Entry) error
	MarkAsError(ctx context.Context, entry *Entry) error
}
