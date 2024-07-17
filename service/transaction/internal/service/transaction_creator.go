package service

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
)

// CreateTransaction defines interface to create transaction.
type CreateTransaction interface {
	// Create creates a new transaction.
	// It needs idempotency key.
	Create(ctx context.Context, transaction entity.Transaction, key string) (string, error)
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
func (tc *TransactionCreator) Create(ctx context.Context, transaction *entity.Transaction, key string) (string, error) {
	if err := tc.validateIdempotencyKey(ctx, key); err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] fail check idempotency key: %s - %v", key, err)
		return "", err
	}

	sanitizeTransaction(transaction)
	if err := validateTransaction(transaction); err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] transaction is invalid: %v", err)
		return "", err
	}

	if err := setTransactionID(ctx, transaction); err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] fail set transaction id: %v", err)
		return "", err
	}
	setTransactionAuditableProperties(transaction)

	err := tc.trxRepo.Insert(ctx, transaction)
	if err != nil {
		app.Logger.Errorf(ctx, "[TransactionCreator-Create] fail save to repository: %v", err)
		return "", err
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
	trx.SenderID = strings.TrimSpace(trx.SenderID)
	trx.ReceiverID = strings.TrimSpace(trx.ReceiverID)
}

func validateTransaction(trx *entity.Transaction) error {
	if trx == nil {
		return entity.ErrEmptyTransaction()
	}
	if trx.SenderID == "" {
		return entity.ErrInvalidSender()
	}
	if trx.ReceiverID == "" {
		return entity.ErrInvalidReceiver()
	}
	if decimal.Zero.Equal(trx.Amount) {
		return entity.ErrInvalidAmount()
	}
	return nil
}

func setTransactionID(ctx context.Context, trx *entity.Transaction) error {
	id, err := generateUniqueID(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[setTransactionID] fail generate unique id: %v", err)
		return entity.ErrInternal("fail to create transaction's ID")
	}
	trx.ID = id
	return nil
}

func generateUniqueID(ctx context.Context) (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		app.Logger.Errorf(ctx, "[generateUniqueID] fail generate ksuid: %v", err)
		return "", entity.ErrInternal("fail to generate unique ID")
	}
	return id.String(), err
}

func setTransactionAuditableProperties(trx *entity.Transaction) {
	trx.CreatedAt = time.Now().UTC()
	trx.UpdatedAt = time.Now().UTC()
	trx.CreatedBy = trx.SenderID
	trx.UpdatedBy = trx.SenderID
}
