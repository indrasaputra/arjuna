package postgres

import (
	"context"

	"github.com/google/uuid"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

// UserOutbox is responsible to connect user_outbox entity with users_outbox table in PostgreSQL.
type UserOutbox struct {
	db uow.DB
}

// NewUserOutbox creates an instance of UserOutbox.
func NewUserOutbox(db uow.DB) *UserOutbox {
	return &UserOutbox{db: db}
}

// InsertWithTx inserts the payload into users_outbox table using transaction.
func (uo *UserOutbox) InsertWithTx(ctx context.Context, tx uow.Tx, payload *entity.UserOutbox) error {
	if tx == nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-InsertWithTx] transaction is not set")
		return entity.ErrInternal("transaction is not set")
	}
	if payload == nil || payload.Payload == nil {
		return entity.ErrEmptyUser()
	}

	query := "INSERT INTO " +
		"users_outbox (id, status, payload, created_at, updated_at, created_by, updated_by) " +
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
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-InsertWithTx] fail insert user with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetAllReady gets all ready records by status in users_outbox table.
// This process uses SELECT FOR UPDATE so be mindful to update the record after using this method.
func (uo *UserOutbox) GetAllReady(ctx context.Context, limit uint) ([]*entity.UserOutbox, error) {
	query := "SELECT id, status, payload FROM users_outbox WHERE status = ? ORDER BY created_at ASC LIMIT ? FOR UPDATE"
	res := []*entity.UserOutbox{}
	err := uo.db.Query(ctx, &res, query, entity.UserOutboxStatusReady, limit)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-GetAllReady] fail get all user's outbox: %v", err)
		return []*entity.UserOutbox{}, entity.ErrInternal(err.Error())
	}
	return res, nil
}

// SetProcessed sets record's status to processed in users_outbox table.
func (uo *UserOutbox) SetProcessed(ctx context.Context, id uuid.UUID) error {
	return uo.SetRecordStatus(ctx, id, entity.UserOutboxStatusProcessed)
}

// SetDelivered sets record's status to delivered in users_outbox table.
func (uo *UserOutbox) SetDelivered(ctx context.Context, id uuid.UUID) error {
	return uo.SetRecordStatus(ctx, id, entity.UserOutboxStatusDelivered)
}

// SetFailed sets record's status to failed in users_outbox table.
func (uo *UserOutbox) SetFailed(ctx context.Context, id uuid.UUID) error {
	return uo.SetRecordStatus(ctx, id, entity.UserOutboxStatusFailed)
}

// SetRecordStatus sets record's status in users_outbox table.
func (uo *UserOutbox) SetRecordStatus(ctx context.Context, id uuid.UUID, status entity.UserOutboxStatus) error {
	query := "UPDATE users_outbox SET status = ? WHERE id = ?"
	_, err := uo.db.Exec(ctx, query, status, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-SetRecordStatus] fail set record's status to %v: %v", status, err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
