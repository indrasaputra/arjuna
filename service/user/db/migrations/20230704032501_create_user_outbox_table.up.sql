BEGIN;

CREATE TYPE status AS ENUM ('READY', 'PROCESSED', 'DELIVERED', 'FAILED');

CREATE TABLE IF NOT EXISTS users_outbox (
    id UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    status STATUS NOT NULL DEFAULT 'READY',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    deleted_by TEXT
);

CREATE INDEX IF NOT EXISTS index_on_users_outbox_on_status_and_created_at ON users_outbox USING btree (
    status, created_at
);

COMMIT;
