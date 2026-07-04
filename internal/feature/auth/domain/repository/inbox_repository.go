package auth

import (
	"context"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/google/uuid"
)

// InboxRepository defines the interface for inbox event persistence (idempotent consumer).
type InboxRepository interface {
	SaveEventIfNotExists(ctx context.Context, event *entity.InboxEvent) (bool, error)
	GetPendingEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error)
	GetRetryableEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error)
	MarkEventProcessing(ctx context.Context, id uuid.UUID) error
	MarkEventProcessed(ctx context.Context, id uuid.UUID) error
	MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error
}
