BEGIN;

DROP INDEX IF EXISTS index_on_users_outbox_on_status_and_created_at
DROP TABLE IF EXISTS users_outbox;
DROP TYPE IF EXISTS status;

COMMIT;
