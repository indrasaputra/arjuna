package repository

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// DeleteUserRepository defines the interface to delete user from repository.
type DeleteUserRepository interface {
	// HardDelete hard-deletes user from repository.
	// If user is not found, it doesn't return error.
	HardDelete(ctx context.Context, id string) error
}

// DeleteUserPostgres defines the interface to delete user from Postgres.
type DeleteUserPostgres interface {
	DeleteUserRepository
	// GetByID gets a user by id.
	GetByID(ctx context.Context, id string) (*entity.User, error)
}

// UserDeleter is responsible to connect user with repositories for delete user purpose.
type UserDeleter struct {
	keycloak DeleteUserRepository
	postgres DeleteUserPostgres
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(keycloak DeleteUserRepository, postgres DeleteUserPostgres) *UserDeleter {
	return &UserDeleter{
		keycloak: keycloak,
		postgres: postgres,
	}
}

// GetByID gets a user by id.
func (ud *UserDeleter) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return ud.postgres.GetByID(ctx, id)
}

// HardDelete hard-deletes a single user from the repository.
// If the user can't be found, it doesn't return error.
func (ud *UserDeleter) HardDelete(ctx context.Context, user *entity.User) error {
	if err := ud.keycloak.HardDelete(ctx, user.KeycloakID); err != nil {
		return err
	}
	return ud.postgres.HardDelete(ctx, user.ID)
}
