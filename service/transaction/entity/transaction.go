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

// TransactionOutboxStatus enumerates transaction outbox status.
type TransactionOutboxStatus string

var (
	// TransactionOutboxStatusReady means ready to be picked up.
	TransactionOutboxStatusReady TransactionOutboxStatus = "READY"
	// TransactionOutboxStatusProcessed means being processed.
	TransactionOutboxStatusProcessed TransactionOutboxStatus = "PROCESSED"
	// TransactionOutboxStatusDelivered means successfully sent to server.
	TransactionOutboxStatusDelivered TransactionOutboxStatus = "DELIVERED"
	// TransactionOutboxStatusFailed means failure.
	TransactionOutboxStatusFailed TransactionOutboxStatus = "FAILED"
)

// TransactionOutbox defines logical data related to transaction outbox.
type TransactionOutbox struct {
	ID      string
	Status  TransactionOutboxStatus
	Payload *Transaction
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
