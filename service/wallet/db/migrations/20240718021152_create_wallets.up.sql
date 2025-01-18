BEGIN;

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    balance NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID,

    CONSTRAINT non_negative_balance CHECK (balance >= 0)
);

CREATE INDEX IF NOT EXISTS index_on_wallets_on_id_and_user_id ON wallets USING btree (
    id, user_id
);

COMMIT;
