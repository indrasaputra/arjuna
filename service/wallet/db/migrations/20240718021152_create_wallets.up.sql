BEGIN;

CREATE TABLE IF NOT EXISTS wallets (
  id            VARCHAR(50)     PRIMARY KEY,
  user_id       VARCHAR(50)     NOT NULL,
  balance       NUMERIC         NOT NULL,
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  created_by    VARCHAR(50),
  updated_by    VARCHAR(50),
  deleted_by    VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS index_on_wallets_on_id_and_user_id ON wallets USING btree (id, user_id);

COMMIT;
