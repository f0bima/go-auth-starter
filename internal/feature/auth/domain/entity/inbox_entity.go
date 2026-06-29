package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// InboxEvent represents a received event to be processed idempotently.
type InboxEvent struct {
	ID            uuid.UUID       `json:"id"`
	MessageID     string          `json:"message_id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
	RetryCount    int             `json:"retry_count"`
	MaxRetries    int             `json:"max_retries"`
	LastError     *string         `json:"last_error,omitempty"`
	NextRetryAt   *time.Time      `json:"next_retry_at,omitempty"`
}

// NewInboxEvent creates a new inbox event.
func NewInboxEvent(messageID string, aggregateID uuid.UUID, aggregateType, eventType string, payload json.RawMessage) *InboxEvent {
	return &InboxEvent{
		ID:            uuid.New(),
		MessageID:     messageID,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		Payload:       payload,
		Status:        "pending",
		MaxRetries:    3,
		CreatedAt:     time.Now(),
	}
}

// MarkProcessing marks the event as being processed.
func (e *InboxEvent) MarkProcessing() {
	e.Status = "processing"
}

// MarkProcessed marks the event as successfully processed.
func (e *InboxEvent) MarkProcessed() {
	now := time.Now()
	e.Status = "processed"
	e.ProcessedAt = &now
}

// MarkFailed marks the event as failed with error.
func (e *InboxEvent) MarkFailed(err error) {
	e.RetryCount++
	errMsg := err.Error()
	e.LastError = &errMsg
	e.Status = "failed"

	if e.RetryCount < e.MaxRetries {
		nextRetry := time.Now().Add(time.Duration(e.RetryCount) * time.Minute)
		e.NextRetryAt = &nextRetry
	}
}

// IsRetryable checks if the event can be retried.
func (e *InboxEvent) IsRetryable() bool {
	return e.Status == "failed" && e.RetryCount < e.MaxRetries
}
