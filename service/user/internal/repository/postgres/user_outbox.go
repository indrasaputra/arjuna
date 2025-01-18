package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
)

// UserOutbox is responsible to connect user_outbox entity with users_outbox table in PostgreSQL.
type UserOutbox struct {
	queries *db.Queries
}

// NewUserOutbox creates an instance of UserOutbox.
func NewUserOutbox(q *db.Queries) *UserOutbox {
	return &UserOutbox{queries: q}
}

// Insert inserts the payload into users_outbox table.
func (uo *UserOutbox) Insert(ctx context.Context, payload *entity.UserOutbox) error {
	if payload == nil || payload.Payload == nil {
		return entity.ErrEmptyUser()
	}

	param := db.CreateUserOutboxParams{
		ID:        payload.ID,
		Status:    db.UserOutboxStatus(payload.Status),
		Payload:   payload.Payload,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
		CreatedBy: payload.CreatedBy,
		UpdatedBy: payload.UpdatedBy,
	}
	err := uo.queries.CreateUserOutbox(ctx, param)
	if uow.IsUniqueViolationError(err) {
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
	param := db.GetAllUserOutboxesForUpdateByStatusParams{
		Status: db.UserOutboxStatus(entity.UserOutboxStatusReady),
		Limit:  int32(limit),
	}
	outboxes, err := uo.queries.GetAllUserOutboxesForUpdateByStatus(ctx, param)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-GetAllReady] fail get all user's outbox: %v", err)
		return []*entity.UserOutbox{}, entity.ErrInternal(err.Error())
	}

	result := make([]*entity.UserOutbox, len(outboxes))
	for i, outbox := range outboxes {
		res := &entity.UserOutbox{
			ID:      outbox.ID,
			Status:  entity.UserOutboxStatus(outbox.Status),
			Payload: outbox.Payload,
		}
		result[i] = res
	}
	return result, nil
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
	param := db.UpdateUserOutboxIDParams{
		ID:     id,
		Status: db.UserOutboxStatus(status),
	}

	err := uo.queries.UpdateUserOutboxID(ctx, param)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUserOutbox-SetRecordStatus] fail set record's status to %v: %v", status, err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
