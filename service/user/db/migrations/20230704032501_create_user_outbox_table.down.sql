BEGIN;

DROP INDEX IF EXISTS index_on_status_and_created_at_on_users_outbox
DROP TABLE IF EXISTS users_outbox;
DROP TYPE IF EXISTS status;

COMMIT;
