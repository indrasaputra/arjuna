// Package uow defines the contract for Unit of Work.
package uow

import (
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is a sentinel error for not found error.
var ErrNotFound = pgx.ErrNoRows
