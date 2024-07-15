BEGIN;

DROP INDEX IF EXISTS index_on_accounts_on_email;
DROP TABLE IF EXISTS accounts;

COMMIT;
