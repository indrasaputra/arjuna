BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    created_by UUID NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_by UUID,

    CONSTRAINT email_length CHECK (LENGTH(email) <= 255),
    CONSTRAINT password_length CHECK (LENGTH(email) <= 255)
);

CREATE INDEX IF NOT EXISTS index_on_accounts_on_email ON accounts USING btree (
    email
);

COMMIT;
