BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
  id            VARCHAR(50)     PRIMARY KEY,
  user_id       VARCHAR(50)     UNIQUE NOT NULL,
  email         VARCHAR(50)     UNIQUE NOT NULL,
  password      VARCHAR(100)    NOT NULL,
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  created_by    VARCHAR(50),
  updated_by    VARCHAR(50),
  deleted_by    VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS index_on_accounts_on_email_and_password ON accounts USING btree (email, password);

COMMIT;
