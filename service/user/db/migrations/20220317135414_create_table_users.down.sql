BEGIN;

DROP INDEX IF EXISTS index_on_email_on_users;
DROP INDEX IF EXISTS index_on_keycloak_id_on_users
DROP TABLE IF EXISTS users;

COMMIT;
