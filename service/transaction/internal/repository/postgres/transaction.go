package postgres

import (
	"context"
	"log/slog"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/db"
)

// Transaction is responsible to connect transaction entity with transactions table in PostgreSQL.
type Transaction struct {
	queries *db.Queries
}

// NewTransaction creates an instance of Transaction.
func NewTransaction(q *db.Queries) *Transaction {
	return &Transaction{queries: q}
}

// Insert inserts a transaction to the database.
func (t *Transaction) Insert(ctx context.Context, trx *entity.Transaction) error {
	if trx == nil {
		return entity.ErrEmptyTransaction()
	}

	param := db.CreateTransactionParams{
		ID:         trx.ID,
		SenderID:   trx.SenderID,
		ReceiverID: trx.ReceiverID,
		Amount:     trx.Amount,
		CreatedAt:  trx.CreatedAt,
		UpdatedAt:  trx.UpdatedAt,
		CreatedBy:  trx.CreatedBy,
		UpdatedBy:  trx.UpdatedBy,
	}
	err := t.queries.CreateTransaction(ctx, param)
	if sdkpostgres.IsUniqueViolationError(err) {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgresTransaction-Insert] fail insert transaction with tx", "error", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// DeleteAll deletes all transactions.
func (t *Transaction) DeleteAll(ctx context.Context) error {
	if err := t.queries.HardDeleteAllTransactions(ctx); err != nil {
		slog.ErrorContext(ctx, "[PostgresTransaction-DeleteAll] fail delete all transactions", "error", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
