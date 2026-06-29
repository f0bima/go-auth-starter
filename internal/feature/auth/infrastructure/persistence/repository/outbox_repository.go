package repository

import (
	"context"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time check
var _ auth.OutboxRepository = (*outboxRepo)(nil)

type outboxRepo struct {
	db *gorm.DB
}

// NewOutboxRepository creates a new Outbox repository.
func NewOutboxRepository(db *gorm.DB) auth.OutboxRepository {
	return &outboxRepo{db: db}
}

func (r *outboxRepo) SaveEvent(ctx context.Context, event *entity.OutboxEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *outboxRepo) GetPendingEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error) {
	var events []*entity.OutboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", "pending", time.Now()).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxRepo) GetRetryableEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error) {
	var events []*entity.OutboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ? AND retry_count < max_retries", "failed", time.Now()).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxRepo) MarkEventProcessed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "processed",
			"processed_at":  now,
			"last_error":    nil,
			"next_retry_at": nil,
		}).Error
}

func (r *outboxRepo) MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var event entity.OutboxEvent
		if err := tx.WithContext(ctx).Where("id = ?", id).First(&event).Error; err != nil {
			return err
		}

		event.MarkFailed(err)

		updates := map[string]interface{}{
			"status":        event.Status,
			"retry_count":   event.RetryCount,
			"last_error":    event.LastError,
			"next_retry_at": event.NextRetryAt,
		}

		return tx.WithContext(ctx).Model(&entity.OutboxEvent{}).Where("id = ?", id).Updates(updates).Error
	})
}
