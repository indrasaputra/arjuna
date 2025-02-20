// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAccount = `-- name: CreateAccount :exec
INSERT INTO accounts (id, user_id, email, password, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateAccountParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	_, err := q.db.Exec(ctx, createAccount,
		arg.ID,
		arg.UserID,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.CreatedBy,
		arg.UpdatedBy,
	)
	return err
}

const getAccountByEmail = `-- name: GetAccountByEmail :one
SELECT id, user_id, email, password, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM accounts WHERE email = $1 LIMIT 1
`

func (q *Queries) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	row := q.db.QueryRow(ctx, getAccountByEmail, email)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return &i, err
}
