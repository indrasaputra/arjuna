CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID,

    CONSTRAINT name_length CHECK (LENGTH(name) <= 100)
);

CREATE INDEX IF NOT EXISTS index_on_users_on_id ON users USING btree (id);

CREATE TYPE user_outbox_status AS ENUM ('READY', 'PROCESSED', 'DELIVERED', 'FAILED');

CREATE TABLE IF NOT EXISTS users_outbox (
    id UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    status USER_OUTBOX_STATUS NOT NULL DEFAULT 'READY',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID
);

CREATE INDEX IF NOT EXISTS index_on_users_outbox_on_status_and_created_at ON users_outbox USING btree (
    status, created_at
);
