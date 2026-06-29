CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    next_retry_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for efficient querying
CREATE INDEX idx_outbox_events_status ON outbox_events(status);
CREATE INDEX idx_outbox_events_next_retry ON outbox_events(next_retry_at) WHERE status = 'failed';
CREATE INDEX idx_outbox_events_aggregate ON outbox_events(aggregate_id, aggregate_type);
CREATE INDEX idx_outbox_events_created ON outbox_events(created_at) WHERE status = 'pending';

-- Partial index for pending events (most common query)
CREATE INDEX idx_outbox_events_pending ON outbox_events(created_at) 
WHERE status = 'pending' AND (next_retry_at IS NULL OR next_retry_at <= NOW());
