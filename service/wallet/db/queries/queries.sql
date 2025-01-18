-- name: CreateWallet :exec
INSERT INTO wallets (id, user_id, balance, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUserWalletForUpdate :one
SELECT * FROM wallets WHERE id = $1 AND user_id = $2 LIMIT 1 FOR NO KEY UPDATE; --noqa

-- name: AddWalletBalance :exec
UPDATE wallets SET balance = balance + @amount WHERE id = $1; --noqa
