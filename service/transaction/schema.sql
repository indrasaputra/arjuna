CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID
);
