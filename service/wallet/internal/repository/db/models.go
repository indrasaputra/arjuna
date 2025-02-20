// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
	Balance   decimal.Decimal
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
