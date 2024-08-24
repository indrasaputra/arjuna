package service

import (
	"context"

	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
)

// DeleteTransaction defines interface to delete transaction.
type DeleteTransaction interface {
	// DeleteAllTransactions deletes all transactions.
	DeleteAllTransactions(ctx context.Context) error
}

// DeleteTransactionRepository defines the interface to delete transaction in repository.
type DeleteTransactionRepository interface {
	// DeleteAll deletes all transactions.
	DeleteAll(ctx context.Context) error
}

// TransactionDeleter is responsible for deleting transaction.
type TransactionDeleter struct {
	repo DeleteTransactionRepository
}

// NewTransactionDeleter creates an instance of TransactionDeleter.
func NewTransactionDeleter(r DeleteTransactionRepository) *TransactionDeleter {
	return &TransactionDeleter{repo: r}
}

// DeleteAllTransactions deletes all transactions.
func (tc *TransactionDeleter) DeleteAllTransactions(ctx context.Context) error {
	err := tc.repo.DeleteAll(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[TransactionDeleter-DeleteAllTransactions] fail delete from repository: %v", err)
		return err
	}
	return nil
}
