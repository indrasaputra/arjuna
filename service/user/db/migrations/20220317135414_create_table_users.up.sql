BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id            VARCHAR(50)     PRIMARY KEY,
  name          VARCHAR(100),
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  created_by    VARCHAR(50),
  updated_by    VARCHAR(50),
  deleted_by    VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS index_on_users_on_id ON users USING btree (id);

COMMIT;
