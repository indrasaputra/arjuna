BEGIN;

DROP INDEX IF EXISTS index_on_wallets_on_id_and_user_id;
DROP TABLE IF EXISTS wallets;

COMMIT;
