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
	// HardDeleteWithTx(ctx context.Context, tx uow.Tx, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
}

// UserDeleter is responsible for deleting a user.
type UserDeleter struct {
	repo      DeleteUserRepository
	txManager uow.TxManager
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(txm uow.TxManager, db DeleteUserRepository) *UserDeleter {
	return &UserDeleter{
		txManager: txm,
		repo:      db,
	}
}

// HardDelete hard-deletes a user in system.
// TODO: Use temporal to call auth service.
func (td *UserDeleter) HardDelete(ctx context.Context, id uuid.UUID) error {
	user, err := td.repo.GetByID(ctx, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserDeleter-HardDelete] fail get user: %v", err)
		return err
	}

	err = td.repo.HardDelete(ctx, user.ID)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserDeleter-HardDelete] fail delete user: %v", err)
	}
	return err
}
