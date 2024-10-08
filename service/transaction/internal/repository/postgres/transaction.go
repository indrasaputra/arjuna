package postgres

import (
	"context"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
)

// Transaction is responsible to connect transaction entity with transactions table in PostgreSQL.
type Transaction struct {
	db uow.DB
}

// NewTransaction creates an instance of Transaction.
func NewTransaction(db uow.DB) *Transaction {
	return &Transaction{db: db}
}

// Insert inserts a transaction to the database.
func (t *Transaction) Insert(ctx context.Context, trx *entity.Transaction) error {
	if trx == nil {
		return entity.ErrEmptyTransaction()
	}

	query := "INSERT INTO " +
		"transactions (id, sender_id, receiver_id, amount, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := t.db.Exec(ctx, query,
		trx.ID,
		trx.SenderID,
		trx.ReceiverID,
		trx.Amount,
		trx.CreatedAt,
		trx.UpdatedAt,
		trx.CreatedBy,
		trx.UpdatedBy,
	)

	if err == sdkpg.ErrAlreadyExist {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresTransaction-Insert] fail insert transaction with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// DeleteAll deletes all transactions.
func (t *Transaction) DeleteAll(ctx context.Context) error {
	query := "DELETE FROM transactions"
	_, err := t.db.Exec(ctx, query)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresTransaction-DeleteAll] fail delete all transactions: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
