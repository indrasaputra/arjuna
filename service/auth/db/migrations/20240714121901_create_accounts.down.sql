BEGIN;

DROP INDEX IF EXISTS index_on_accounts_on_email_and_password;
DROP TABLE IF EXISTS accounts;

COMMIT;
