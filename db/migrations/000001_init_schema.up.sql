CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Watermill Explicit Schema (Production-Grade)
CREATE TABLE IF NOT EXISTS "watermill_events_to_forward" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_events_to_forward" (
    consumer_group VARCHAR(255) NOT NULL,
    offset_acked BIGINT,
    last_processed_transaction_id xid8 NOT NULL,
    PRIMARY KEY(consumer_group)
);
