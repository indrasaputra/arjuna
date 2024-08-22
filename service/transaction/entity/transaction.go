package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction defines logical data related to transaction.
type Transaction struct {
	Auditable
	Amount     decimal.Decimal
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

// Auditable defines logical data related to audit.
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}
