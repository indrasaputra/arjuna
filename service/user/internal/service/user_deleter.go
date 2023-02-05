package service

import (
	"context"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

// DeleteUser defines the interface to delete a user.
type DeleteUser interface {
	// HardDelete hard-deletes all users.
	HardDelete(ctx context.Context, id string) error
}

// DeleteUserProvider defines the interface to delete user from the provider.
type DeleteUserProvider interface {
	// HardDelete hard-deletes a single user from the provider.
	// If the user can't be found, it doesn't return error.
	HardDelete(ctx context.Context, id string) error
}

// DeleteUserRepository defines interface to delete user from repository.
type DeleteUserRepository interface {
	// GetByID gets a user by ID.
	GetByID(ctx context.Context, id string) (*entity.User, error)
	// HardDelete hard-deletes a single user from the repository.
	// If the user can't be found, it doesn't return error.
	HardDelete(ctx context.Context, tx uow.Tx, id string) error
}

// UserDeleter is responsible for deleting a user.
type UserDeleter struct {
	database DeleteUserRepository
	keycloak DeleteUserProvider
	unit     uow.UnitOfWork
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(unit uow.UnitOfWork, db DeleteUserRepository, kc DeleteUserProvider) *UserDeleter {
	return &UserDeleter{
		database: db,
		keycloak: kc,
		unit:     unit,
	}
}

// HardDelete hard-deletes a user in system.
func (td *UserDeleter) HardDelete(ctx context.Context, id string) error {
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
	if err := td.database.HardDelete(ctx, tx, user.ID); err != nil {
		return err
	}
	return td.keycloak.HardDelete(ctx, user.KeycloakID)
}
