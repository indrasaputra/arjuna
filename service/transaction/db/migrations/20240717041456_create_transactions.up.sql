BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
  id            VARCHAR(50)     PRIMARY KEY,
  sender_id     VARCHAR(50)     NOT NULL,
  receiver_id   VARCHAR(50)     NOT NULL,
  amount        NUMERIC         NOT NULL,
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  created_by    VARCHAR(50),
  updated_by    VARCHAR(50),
  deleted_by    VARCHAR(50)
);

COMMIT;
