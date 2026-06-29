package repository

import (
	"context"
	"errors"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time check
var _ auth.InboxRepository = (*inboxRepo)(nil)

type inboxRepo struct {
	db *gorm.DB
}

// NewInboxRepository creates a new Inbox repository.
func NewInboxRepository(db *gorm.DB) auth.InboxRepository {
	return &inboxRepo{db: db}
}

// SaveEventIfNotExists saves the event only if it doesn't exist (idempotent).
// Returns (isNew, error).
func (r *inboxRepo) SaveEventIfNotExists(ctx context.Context, event *entity.InboxEvent) (bool, error) {
	// Check if message_id already exists
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.InboxEvent{}).
		Where("message_id = ?", event.MessageID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		// Event already exists (duplicate)
		return false, nil
	}

	// Insert new event
	err = r.db.WithContext(ctx).Create(event).Error
	if err != nil {
		// Check for unique constraint violation (race condition)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *inboxRepo) GetPendingEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error) {
	var events []*entity.InboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", "pending", time.Now()).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *inboxRepo) GetRetryableEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error) {
	var events []*entity.InboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ? AND retry_count < max_retries", "failed", time.Now()).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *inboxRepo) MarkEventProcessing(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.InboxEvent{}).
		Where("id = ? AND status = ?", id, "pending").
		Update("status", "processing").Error
}

func (r *inboxRepo) MarkEventProcessed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.InboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "processed",
			"processed_at":  now,
			"last_error":    nil,
			"next_retry_at": nil,
		}).Error
}

func (r *inboxRepo) MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var event entity.InboxEvent
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

		return tx.WithContext(ctx).Model(&entity.InboxEvent{}).Where("id = ?", id).Updates(updates).Error
	})
}
