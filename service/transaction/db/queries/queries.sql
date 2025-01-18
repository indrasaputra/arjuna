-- name: CreateTransaction :exec
INSERT INTO transactions (id, sender_id, receiver_id, amount, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: HardDeleteAllTransactions :exec
DELETE FROM transactions;
