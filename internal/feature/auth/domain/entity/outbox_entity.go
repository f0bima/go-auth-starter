package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OutboxEvent represents an event to be published asynchronously.
type OutboxEvent struct {
	ID            uuid.UUID       `json:"id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Status        string          `json:"status"` // pending, processing, processed, failed
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
	RetryCount    int             `json:"retry_count"`
	MaxRetries    int             `json:"max_retries"`
	LastError     *string         `json:"last_error,omitempty"`
	NextRetryAt   *time.Time      `json:"next_retry_at,omitempty"`
}

// NewOutboxEvent creates a new outbox event.
func NewOutboxEvent(aggregateID uuid.UUID, aggregateType, eventType string, payload interface{}) (*OutboxEvent, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &OutboxEvent{
		ID:            uuid.New(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		Payload:       payloadBytes,
		Status:        "pending",
		MaxRetries:    3,
		CreatedAt:     time.Now(),
	}, nil
}

// MarkProcessed marks the event as processed.
func (e *OutboxEvent) MarkProcessed() {
	now := time.Now()
	e.Status = "processed"
	e.ProcessedAt = &now
}

// MarkFailed marks the event as failed with error.
func (e *OutboxEvent) MarkFailed(err error) {
	e.RetryCount++
	errMsg := err.Error()
	e.LastError = &errMsg
	e.Status = "failed"

	if e.RetryCount < e.MaxRetries {
		// Schedule retry with exponential backoff
		nextRetry := time.Now().Add(time.Duration(e.RetryCount) * time.Minute)
		e.NextRetryAt = &nextRetry
	}
}

// IsRetryable checks if the event can be retried.
func (e *OutboxEvent) IsRetryable() bool {
	return e.Status == "failed" && e.RetryCount < e.MaxRetries
}
