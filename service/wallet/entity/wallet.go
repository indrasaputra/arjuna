package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Wallet defines logical data related to wallet.
type Wallet struct {
	Auditable
	Balance decimal.Decimal `json:"balance"`
	ID      uuid.UUID       `json:"id"`
	UserID  uuid.UUID       `json:"user_id"`
}

// TopupWallet defines logical data related to topup wallet.
type TopupWallet struct {
	Amount         decimal.Decimal
	IdempotencyKey string
	WalletID       uuid.UUID
	UserID         uuid.UUID
}

// TransferWallet defines logical data related to transfer wallet.
type TransferWallet struct {
	Amount           decimal.Decimal
	SenderID         uuid.UUID
	SenderWalletID   uuid.UUID
	ReceiverID       uuid.UUID
	ReceiverWalletID uuid.UUID
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
