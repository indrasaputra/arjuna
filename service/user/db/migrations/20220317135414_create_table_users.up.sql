BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id            VARCHAR(50)     PRIMARY KEY,
  name          VARCHAR(100),
  email         VARCHAR(100)    UNIQUE NOT NULL,
  password      TEXT,
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  created_by    VARCHAR(50),
  updated_by    VARCHAR(50),
  deleted_by    VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS index_on_email_on_users ON users USING btree (email);

COMMIT;
