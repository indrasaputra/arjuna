package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

// DeleteUser defines the interface to delete a user.
type DeleteUser interface {
	// HardDelete hard-deletes all users.
	HardDelete(ctx context.Context, id uuid.UUID) error
}

// DeleteUserRepository defines interface to delete user from repository.
type DeleteUserRepository interface {
	// GetByID gets a user by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	// HardDeleteWithTx hard-deletes a single user from the repository using transaction.
	// If the user can't be found, it doesn't return error.
	HardDeleteWithTx(ctx context.Context, tx uow.Tx, id uuid.UUID) error
}

// UserDeleter is responsible for deleting a user.
type UserDeleter struct {
	database DeleteUserRepository
	unit     uow.UnitOfWork
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(unit uow.UnitOfWork, db DeleteUserRepository) *UserDeleter {
	return &UserDeleter{
		database: db,
		unit:     unit,
	}
}

// HardDelete hard-deletes a user in system.
// TODO: Use temporal to call auth service.
func (td *UserDeleter) HardDelete(ctx context.Context, id uuid.UUID) error {
	user, err := td.database.GetByID(ctx, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserDeleter-HardDelete] fail get user: %v", err)
		return err
	}

	tx, err := td.unit.Begin(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserDeleter-HardDelete] fail init transaction: %v", err)
		return err
	}

	err = td.hardDelete(ctx, tx, user)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserDeleter-HardDelete] fail delete user: %v", err)
	}
	return td.unit.Finish(ctx, tx, err)
}

func (td *UserDeleter) hardDelete(ctx context.Context, tx uow.Tx, user *entity.User) error {
	return td.database.HardDeleteWithTx(ctx, tx, user.ID)
}
