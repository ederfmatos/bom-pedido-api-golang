package outbox

import (
	"context"
	"time"
)

type (
	Entry struct {
		Id              string     `json:"id,omitempty" bson:"_id"`
		Name            string     `json:"name,omitempty" bson:"name"`
		Data            string     `json:"data,omitempty" bson:"data"`
		CreatedAt       time.Time  `json:"createdAt,omitempty" bson:"createdAt"`
		Status          string     `json:"status,omitempty" bson:"status"`
		ProcessedAt     *time.Time `json:"processedAt" bson:"processedAt"`
		LastAttemptTime *time.Time `json:"lastAttemptTime" bson:"lastAttemptTime"`
	}

	Repository interface {
		Save(ctx context.Context, entry *Entry) error
		Get(ctx context.Context, id string) (*Entry, error)
		Update(ctx context.Context, entry *Entry) error
	}
)

func (entry *Entry) MarkAsError() {
	now := time.Now()
	entry.LastAttemptTime = &now
	entry.Status = "ERROR"
}

func (entry *Entry) MarkAsProcessed() {
	now := time.Now()
	entry.ProcessedAt = &now
	entry.Status = "PROCESSED"
}
