package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction defines logical data related to transaction.
type Transaction struct {
	ID         string
	SenderID   string
	ReceiverID string
	Amount     decimal.Decimal
	Auditable
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
