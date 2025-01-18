BEGIN;

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

COMMIT;
