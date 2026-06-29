package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
)

// EventPublisher defines the interface for publishing events to message broker.
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

// OutboxProcessor polls the outbox table and publishes pending events.
type OutboxProcessor struct {
	outboxRepo   auth.OutboxRepository
	publisher    EventPublisher
	batchSize    int
	pollInterval time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	isRunning    bool
	mu           sync.Mutex
}

// NewOutboxProcessor creates a new outbox processor.
func NewOutboxProcessor(
	outboxRepo auth.OutboxRepository,
	publisher EventPublisher,
	batchSize int,
	pollInterval time.Duration,
) *OutboxProcessor {
	if batchSize == 0 {
		batchSize = 10
	}
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	return &OutboxProcessor{
		outboxRepo:   outboxRepo,
		publisher:    publisher,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
	}
}

// Start begins polling the outbox table.
func (p *OutboxProcessor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("outbox processor already running")
	}
	p.isRunning = true
	p.mu.Unlock()

	slog.Info("Outbox processor started", "poll_interval", p.pollInterval, "batch_size", p.batchSize)

	p.wg.Add(1)
	go p.run(ctx)

	return nil
}

// Stop gracefully stops the outbox processor.
func (p *OutboxProcessor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return
	}

	slog.Info("Stopping outbox processor...")
	close(p.stopChan)
	p.wg.Wait()
	p.isRunning = false
	slog.Info("Outbox processor stopped")
}

func (p *OutboxProcessor) run(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, stopping outbox processor")
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.processOutbox(ctx)
		}
	}
}

func (p *OutboxProcessor) processOutbox(ctx context.Context) {
	// Get pending events
	events, err := p.outboxRepo.GetPendingEvents(ctx, p.batchSize)
	if err != nil {
		slog.Error("Failed to get pending events", "error", err)
		return
	}

	if len(events) == 0 {
		return
	}

	slog.Info("Processing outbox events", "count", len(events))

	for _, event := range events {
		if err := p.publishEvent(ctx, event); err != nil {
			slog.Error("Failed to publish event",
				"event_id", event.ID,
				"event_type", event.EventType,
				"error", err,
			)

			// Mark as failed
			if markErr := p.outboxRepo.MarkEventFailed(ctx, event.ID, err); markErr != nil {
				slog.Error("Failed to mark event as failed",
					"event_id", event.ID,
					"error", markErr,
				)
			}
		}
	}
}

func (p *OutboxProcessor) publishEvent(ctx context.Context, event *entity.OutboxEvent) error {
	// Parse payload
	var payload map[string]interface{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	// Add metadata
	payload["event_id"] = event.ID.String()
	payload["aggregate_id"] = event.AggregateID.String()
	payload["aggregate_type"] = event.AggregateType
	payload["event_type"] = event.EventType
	payload["occurred_at"] = event.CreatedAt

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload with metadata: %w", err)
	}

	// Publish to message broker
	if err := p.publisher.Publish(ctx, event.EventType, payloadBytes); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	// Mark as processed
	if err := p.outboxRepo.MarkEventProcessed(ctx, event.ID); err != nil {
		return fmt.Errorf("mark event processed: %w", err)
	}

	slog.Info("Event published successfully",
		"event_id", event.ID,
		"event_type", event.EventType,
	)

	return nil
}
