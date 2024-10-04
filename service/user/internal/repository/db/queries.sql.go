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

const createUser = `-- name: CreateUser :exec
INSERT INTO
users (id, name, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateUserParams struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.CreatedBy,
		arg.UpdatedBy,
	)
	return err
}

const createUserOutbox = `-- name: CreateUserOutbox :exec
INSERT INTO
users_outbox (id, status, payload, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateUserOutboxParams struct {
	ID        uuid.UUID
	Status    Status
	Payload   []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}

func (q *Queries) CreateUserOutbox(ctx context.Context, arg CreateUserOutboxParams) error {
	_, err := q.db.Exec(ctx, createUserOutbox,
		arg.ID,
		arg.Status,
		arg.Payload,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.CreatedBy,
		arg.UpdatedBy,
	)
	return err
}

const getAllUserOutboxesForUpdateByStatus = `-- name: GetAllUserOutboxesForUpdateByStatus :many
SELECT id, payload, status, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users_outbox WHERE status = $1 FOR UPDATE
`

func (q *Queries) GetAllUserOutboxesForUpdateByStatus(ctx context.Context, status Status) ([]*UsersOutbox, error) {
	rows, err := q.db.Query(ctx, getAllUserOutboxesForUpdateByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*UsersOutbox
	for rows.Next() {
		var i UsersOutbox
		if err := rows.Scan(
			&i.ID,
			&i.Payload,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users LIMIT $1
`

func (q *Queries) GetAllUsers(ctx context.Context, limit int32) ([]*User, error) {
	rows, err := q.db.Query(ctx, getAllUsers, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return &i, err
}

const hardDeleteUserByID = `-- name: HardDeleteUserByID :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) HardDeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, hardDeleteUserByID, id)
	return err
}

const updateUserOutboxID = `-- name: UpdateUserOutboxID :exec
UPDATE users_outbox SET status = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3
`

type UpdateUserOutboxIDParams struct {
	Status    Status
	UpdatedBy string
	ID        uuid.UUID
}

func (q *Queries) UpdateUserOutboxID(ctx context.Context, arg UpdateUserOutboxIDParams) error {
	_, err := q.db.Exec(ctx, updateUserOutboxID, arg.Status, arg.UpdatedBy, arg.ID)
	return err
}
