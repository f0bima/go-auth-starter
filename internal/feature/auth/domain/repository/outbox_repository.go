package auth

import (
	"context"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/google/uuid"
)

// OutboxRepository defines the interface for outbox event persistence.
type OutboxRepository interface {
	SaveEvent(ctx context.Context, event *entity.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	GetRetryableEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, id uuid.UUID) error
	MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error
}
