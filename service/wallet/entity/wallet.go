package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Wallet defines logical data related to wallet.
type Wallet struct {
	ID      string
	UserID  string
	Balance decimal.Decimal
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
