// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	DeletedBy  *uuid.UUID
	Amount     decimal.Decimal
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	CreatedBy  uuid.UUID
	UpdatedBy  uuid.UUID
}
