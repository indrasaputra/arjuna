package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
)

// CreateTransaction defines interface to create transaction.
type CreateTransaction interface {
	// Create creates a new transaction.
	// It needs idempotency key.
	Create(ctx context.Context, transaction *entity.Transaction, key string) (uuid.UUID, error)
}

// IdempotencyKeyRepository defines  interface for idempotency check flow and repository.
type IdempotencyKeyRepository interface {
	// Exists check if given key exists in repository.
	Exists(ctx context.Context, key string) (bool, error)
}

// CreateTransactionRepository defines the interface to insert transaction to repository.
type CreateTransactionRepository interface {
	// Insert inserts a transaction.
	Insert(ctx context.Context, transaction *entity.Transaction) error
}

// TransactionCreator is responsible for creating a new transaction.
type TransactionCreator struct {
	trxRepo CreateTransactionRepository
	keyRepo IdempotencyKeyRepository
}

// NewTransactionCreator creates an instance of TransactionCreator.
func NewTransactionCreator(t CreateTransactionRepository, k IdempotencyKeyRepository) *TransactionCreator {
	return &TransactionCreator{trxRepo: t, keyRepo: k}
}

// Create creates a new transaction.
// It needs idempotency key.
func (tc *TransactionCreator) Create(ctx context.Context, transaction *entity.Transaction, key string) (uuid.UUID, error) {
	if err := tc.validateIdempotencyKey(ctx, key); err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] fail check idempotency key: %s - %v", key, err)
		return uuid.Nil, err
	}

	sanitizeTransaction(transaction)
	if err := validateTransaction(transaction); err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] transaction is invalid: %v", err)
		return uuid.Nil, err
	}

	setTransactionID(transaction)
	setTransactionAuditableProperties(transaction)

	err := tc.trxRepo.Insert(ctx, transaction)
	if err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] fail save to repository: %v", err)
		return uuid.Nil, err
	}
	return transaction.ID, nil
}

func (tc *TransactionCreator) validateIdempotencyKey(ctx context.Context, key string) error {
	res, err := tc.keyRepo.Exists(ctx, key)
	if err != nil {
		return err
	}
	if res {
		return entity.ErrAlreadyExists()
	}
	return nil
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
	trx.CreatedBy = trx.SenderID.String()
	trx.UpdatedBy = trx.SenderID.String()
}
