-- name: CreateUser :exec
INSERT INTO
users (id, name, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users LIMIT $1;

-- name: HardDeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: CreateUserOutbox :exec
INSERT INTO
users_outbox (id, status, payload, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetAllUserOutboxesForUpdateByStatus :many
SELECT * FROM users_outbox WHERE status = $1 ORDER BY created_at ASC LIMIT $2;

-- name: UpdateUserOutboxID :exec
UPDATE users_outbox SET status = $1, updated_at = NOW() WHERE id = $2;
