BEGIN;

DROP INDEX IF EXISTS index_on_transactions_outbox_on_status_and_created_at
DROP TABLE IF EXISTS transactions_outbox;
DROP TYPE IF EXISTS status;

COMMIT;
