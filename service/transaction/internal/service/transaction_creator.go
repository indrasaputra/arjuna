package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/service/transaction/entity"
)

// CreateTransaction defines interface to create transaction.
type CreateTransaction interface {
	// Create creates a new transaction.
	Create(ctx context.Context, transaction *entity.Transaction) (uuid.UUID, error)
}

// CreateTransactionRepository defines the interface to insert transaction to repository.
type CreateTransactionRepository interface {
	// Insert inserts a transaction.
	Insert(ctx context.Context, transaction *entity.Transaction) error
}

// TransactionCreator is responsible for creating a new transaction.
type TransactionCreator struct {
	trxRepo CreateTransactionRepository
}

// NewTransactionCreator creates an instance of TransactionCreator.
func NewTransactionCreator(t CreateTransactionRepository) *TransactionCreator {
	return &TransactionCreator{trxRepo: t}
}

// Create creates a new transaction.
func (tc *TransactionCreator) Create(ctx context.Context, transaction *entity.Transaction) (uuid.UUID, error) {
	sanitizeTransaction(transaction)
	if err := validateTransaction(transaction); err != nil {
		slog.ErrorContext(ctx, "[TransactionCreator-Create] transaction is invalid", "error", err)
		return uuid.Nil, err
	}

	setTransactionID(transaction)
	setTransactionAuditableProperties(transaction)

	err := tc.trxRepo.Insert(ctx, transaction)
	if err != nil {
		slog.ErrorContext(ctx, "[TransactionCreator-Create] fail save to repository", "error", err)
		return uuid.Nil, err
	}
	return transaction.ID, nil
}

func sanitizeTransaction(trx *entity.Transaction) {
	if trx == nil {
		return
	}
}

func validateTransaction(trx *entity.Transaction) error {
	if trx == nil {
		return entity.ErrEmptyTransaction()
	}
	if trx.SenderID == uuid.Nil {
		return entity.ErrInvalidSender()
	}
	if trx.ReceiverID == uuid.Nil {
		return entity.ErrInvalidReceiver()
	}
	if decimal.Zero.Equal(trx.Amount) {
		return entity.ErrInvalidAmount()
	}
	return nil
}

func setTransactionID(trx *entity.Transaction) {
	trx.ID = generateUniqueID()
}

func generateUniqueID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func setTransactionAuditableProperties(trx *entity.Transaction) {
	trx.CreatedAt = time.Now().UTC()
	trx.UpdatedAt = time.Now().UTC()
	trx.CreatedBy = trx.SenderID
	trx.UpdatedBy = trx.SenderID
}
