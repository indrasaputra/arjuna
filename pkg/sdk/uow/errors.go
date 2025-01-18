// Package uow defines the contract for Unit of Work.
package uow

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrNotFound is a sentinel error for not found error.
var ErrNotFound = pgx.ErrNoRows

// IsUniqueViolationError checks if the error is a unique violation error.
func IsUniqueViolationError(err error) bool {
	if err == nil {
		return false
	}
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}
