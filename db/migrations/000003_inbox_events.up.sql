CREATE TABLE IF NOT EXISTS inbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id VARCHAR(255) UNIQUE NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    last_error TEXT,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT chk_status CHECK (status IN ('pending', 'processing', 'processed', 'failed'))
);

-- Unique constraint on message_id ensures idempotency
CREATE UNIQUE INDEX idx_inbox_events_message_id ON inbox_events(message_id);

-- Indexes for efficient querying
CREATE INDEX idx_inbox_events_status ON inbox_events(status);
CREATE INDEX idx_inbox_events_next_retry ON inbox_events(next_retry_at) WHERE status = 'failed';
CREATE INDEX idx_inbox_events_aggregate ON inbox_events(aggregate_id, aggregate_type);
CREATE INDEX idx_inbox_events_type ON inbox_events(event_type) WHERE status = 'pending';
