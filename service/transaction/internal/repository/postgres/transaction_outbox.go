package postgres

import (
	"context"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
)

// TransactionOutbox is responsible to connect transaction_outbox entity with transactions_outbox table in PostgreSQL.
type TransactionOutbox struct {
	db uow.DB
}

// NewTransactionOutbox creates an instance of TransactionOutbox.
func NewTransactionOutbox(db uow.DB) *TransactionOutbox {
	return &TransactionOutbox{db: db}
}

// InsertWithTx inserts the payload into transactions_outbox table using transaction.
func (to *TransactionOutbox) InsertWithTx(ctx context.Context, tx uow.Tx, payload *entity.TransactionOutbox) error {
	if tx == nil {
		app.Logger.Errorf(ctx, "[PostgresTransactionOutbox-InsertWithTx] transaction is not set")
		return entity.ErrInternal("transaction is not set")
	}
	if payload == nil || payload.Payload == nil {
		return entity.ErrEmptyTransaction()
	}

	query := "INSERT INTO " +
		"transactions_outbox (id, status, payload, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := tx.Exec(ctx, query,
		payload.ID,
		payload.Status,
		payload.Payload,
		payload.CreatedAt,
		payload.UpdatedAt,
		payload.CreatedBy,
		payload.UpdatedBy,
	)

	if err == sdkpg.ErrAlreadyExist {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresTransactionOutbox-InsertWithTx] fail insert transaction with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetAllReady gets all ready records by status in transactions_outbox table.
// This process uses SELECT FOR UPDATE so be mindful to update the record after using this method.
func (to *TransactionOutbox) GetAllReady(ctx context.Context, limit uint) ([]*entity.TransactionOutbox, error) {
	query := "SELECT id, status, payload FROM transactions_outbox WHERE status = ? ORDER BY created_at ASC LIMIT ? FOR UPDATE"
	res := []*entity.TransactionOutbox{}
	err := to.db.Query(ctx, &res, query, entity.TransactionOutboxStatusReady, limit)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresTransactionOutbox-GetAllReady] fail get all transaction's outbox: %v", err)
		return []*entity.TransactionOutbox{}, entity.ErrInternal(err.Error())
	}
	return res, nil
}

// SetProcessed sets record's status to processed in transactions_outbox table.
func (to *TransactionOutbox) SetProcessed(ctx context.Context, id string) error {
	return to.SetRecordStatus(ctx, id, entity.TransactionOutboxStatusProcessed)
}

// SetDelivered sets record's status to delivered in transactions_outbox table.
func (to *TransactionOutbox) SetDelivered(ctx context.Context, id string) error {
	return to.SetRecordStatus(ctx, id, entity.TransactionOutboxStatusDelivered)
}

// SetFailed sets record's status to failed in transactions_outbox table.
func (to *TransactionOutbox) SetFailed(ctx context.Context, id string) error {
	return to.SetRecordStatus(ctx, id, entity.TransactionOutboxStatusFailed)
}

// SetRecordStatus sets record's status in transactions_outbox table.
func (to *TransactionOutbox) SetRecordStatus(ctx context.Context, id string, status entity.TransactionOutboxStatus) error {
	query := "UPDATE transactions_outbox SET status = ? WHERE id = ?"
	_, err := to.db.Exec(ctx, query, status, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresTransactionOutbox-SetRecordStatus] fail set record's status to %v: %v", status, err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
