BEGIN;

CREATE TYPE status AS ENUM ('READY', 'PROCESSED', 'DELIVERED', 'FAILED');

CREATE TABLE IF NOT EXISTS transactions_outbox (
  id                VARCHAR(50)     PRIMARY KEY,
  payload           JSONB           NOT NULL,
  status            status          NOT NULL    DEFAULT 'READY',
  created_at        TIMESTAMP,
  updated_at        TIMESTAMP,
  deleted_at        TIMESTAMP,
  created_by        VARCHAR(50),
  updated_by        VARCHAR(50),
  deleted_by        VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS index_on_transactions_outbox_on_status_and_created_at ON transactions_outbox USING btree (status, created_at);

COMMIT;
