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

// EventHandler defines the interface for handling specific event types.
type EventHandler interface {
	Handle(ctx context.Context, event *entity.InboxEvent) error
	CanHandle(eventType string) bool
}

// InboxProcessor polls the inbox table and processes pending events idempotently.
type InboxProcessor struct {
	inboxRepo    auth.InboxRepository
	handlers     []EventHandler
	batchSize    int
	pollInterval time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	isRunning    bool
	mu           sync.Mutex
}

// NewInboxProcessor creates a new inbox processor.
func NewInboxProcessor(
	inboxRepo auth.InboxRepository,
	handlers []EventHandler,
	batchSize int,
	pollInterval time.Duration,
) *InboxProcessor {
	if batchSize == 0 {
		batchSize = 10
	}
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	return &InboxProcessor{
		inboxRepo:    inboxRepo,
		handlers:     handlers,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
	}
}

// AddHandler registers an event handler.
func (p *InboxProcessor) AddHandler(handler EventHandler) {
	p.handlers = append(p.handlers, handler)
}

// Start begins polling the inbox table.
func (p *InboxProcessor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("inbox processor already running")
	}
	p.isRunning = true
	p.mu.Unlock()

	slog.Info("Inbox processor started", "poll_interval", p.pollInterval, "batch_size", p.batchSize, "handlers", len(p.handlers))

	p.wg.Add(1)
	go p.run(ctx)

	return nil
}

// Stop gracefully stops the inbox processor.
func (p *InboxProcessor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return
	}

	slog.Info("Stopping inbox processor...")
	close(p.stopChan)
	p.wg.Wait()
	p.isRunning = false
	slog.Info("Inbox processor stopped")
}

func (p *InboxProcessor) run(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, stopping inbox processor")
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.processInbox(ctx)
		}
	}
}

func (p *InboxProcessor) processInbox(ctx context.Context) {
	// Get pending events
	events, err := p.inboxRepo.GetPendingEvents(ctx, p.batchSize)
	if err != nil {
		slog.Error("Failed to get pending inbox events", "error", err)
		return
	}

	if len(events) == 0 {
		return
	}

	slog.Info("Processing inbox events", "count", len(events))

	for _, event := range events {
		if err := p.processEvent(ctx, event); err != nil {
			slog.Error("Failed to process inbox event",
				"event_id", event.ID,
				"message_id", event.MessageID,
				"event_type", event.EventType,
				"error", err,
			)

			// Mark as failed
			if markErr := p.inboxRepo.MarkEventFailed(ctx, event.ID, err); markErr != nil {
				slog.Error("Failed to mark inbox event as failed",
					"event_id", event.ID,
					"error", markErr,
				)
			}
		}
	}
}

func (p *InboxProcessor) processEvent(ctx context.Context, event *entity.InboxEvent) error {
	// Find handler for this event type
	var handler EventHandler
	for _, h := range p.handlers {
		if h.CanHandle(event.EventType) {
			handler = h
			break
		}
	}

	if handler == nil {
		// No handler found - mark as processed to avoid infinite retry
		slog.Warn("No handler found for event type, marking as processed",
			"event_type", event.EventType,
		)
		return p.inboxRepo.MarkEventProcessed(ctx, event.ID)
	}

	// Mark as processing
	if err := p.inboxRepo.MarkEventProcessing(ctx, event.ID); err != nil {
		return fmt.Errorf("mark event processing: %w", err)
	}

	// Handle the event
	if err := handler.Handle(ctx, event); err != nil {
		return fmt.Errorf("handle event: %w", err)
	}

	// Mark as processed
	if err := p.inboxRepo.MarkEventProcessed(ctx, event.ID); err != nil {
		return fmt.Errorf("mark event processed: %w", err)
	}

	slog.Info("Inbox event processed successfully",
		"event_id", event.ID,
		"message_id", event.MessageID,
		"event_type", event.EventType,
	)

	return nil
}

// ParsePayload unmarshals the event payload into the provided struct.
func ParsePayload(event *entity.InboxEvent, target interface{}) error {
	return json.Unmarshal(event.Payload, target)
}
