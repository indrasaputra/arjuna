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
	db     uow.Tr
	getter uow.TxGetter
}

// NewUserOutbox creates an instance of UserOutbox.
func NewUserOutbox(db uow.Tr, g uow.TxGetter) *UserOutbox {
	return &UserOutbox{db: db, getter: g}
}

// Insert inserts the payload into users_outbox table.
func (uo *UserOutbox) Insert(ctx context.Context, payload *entity.UserOutbox) error {
	if payload == nil || payload.Payload == nil {
		return entity.ErrEmptyUser()
	}

	tx := uo.getter.DefaultTrOrDB(ctx, uo.db)
	query := "INSERT INTO " +
		"users_outbox (id, status, payload, created_at, updated_at, created_by, updated_by) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"

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
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-Insert] fail insert user with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetAllReady gets all ready records by status in users_outbox table.
// This process uses SELECT FOR UPDATE so be mindful to update the record after using this method.
func (uo *UserOutbox) GetAllReady(ctx context.Context, limit uint) ([]*entity.UserOutbox, error) {
	tx := uo.getter.DefaultTrOrDB(ctx, uo.db)

	query := "SELECT id, status, payload FROM users_outbox WHERE status = $1 ORDER BY created_at ASC LIMIT $2 FOR UPDATE"
	res := []*entity.UserOutbox{}
	rows, err := tx.Query(ctx, query, entity.UserOutboxStatusReady, limit)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-GetAllReady] fail get all user's outbox: %v", err)
		return []*entity.UserOutbox{}, entity.ErrInternal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var tmp entity.UserOutbox
		if err := rows.Scan(&tmp.ID, &tmp.Status, &tmp.Payload); err != nil {
			app.Logger.Errorf(ctx, "[PostgresUserOutbox-GetAll] scan rows error: %v", err)
			return []*entity.UserOutbox{}, entity.ErrInternal(err.Error())
		}
		res = append(res, &tmp)
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
	tx := uo.getter.DefaultTrOrDB(ctx, uo.db)
	query := "UPDATE users_outbox SET status = $1 WHERE id = $2"
	_, err := tx.Exec(ctx, query, status, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-SetRecordStatus] fail set record's status to %v: %v", status, err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
