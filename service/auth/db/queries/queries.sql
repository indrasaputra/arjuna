-- name: CreateAccount :exec
INSERT INTO accounts (id, user_id, email, password, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetAccountByEmail :one
SELECT * FROM accounts WHERE email = $1 LIMIT 1;
