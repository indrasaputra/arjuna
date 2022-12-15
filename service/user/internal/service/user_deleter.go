package service

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// DeleteUser defines the interface to delete a user.
type DeleteUser interface {
	// HardDelete hard-deletes all users.
	HardDelete(ctx context.Context, id string) error
}

// DeleteUserRepository defines the interface to delete user from the repository.
type DeleteUserRepository interface {
	// HardDelete hard-deletes a single user from the repository.
	// If the user can't be found, it doesn't return error.
	HardDelete(ctx context.Context, id string) error
}

// DeleteUserDatabase defines interface to delete user from database.
type DeleteUserDatabase interface {
	// GetByID gets a user by ID.
	GetByID(ctx context.Context, id string) (*entity.User, error)
	DeleteUserRepository
}

// Transactor defines the interface for database transaction.
type Transactor interface {
	// WithinTransaction executes the fn atomically.
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// UserDeleter is responsible for deleting a user.
type UserDeleter struct {
	database   DeleteUserDatabase
	keycloak   DeleteUserRepository
	transactor Transactor
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(db DeleteUserDatabase, kc DeleteUserRepository, tx Transactor) *UserDeleter {
	return &UserDeleter{
		database:   db,
		keycloak:   kc,
		transactor: tx,
	}
}

// HardDelete hard-deletes a user in system.
func (td *UserDeleter) HardDelete(ctx context.Context, id string) error {
	user, err := td.database.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return td.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		return td.hardDelete(txCtx, user)
	})
}

func (td *UserDeleter) hardDelete(ctx context.Context, user *entity.User) error {
	if err := td.database.HardDelete(ctx, user.ID); err != nil {
		return err
	}
	return td.keycloak.HardDelete(ctx, user.KeycloakID)
}
