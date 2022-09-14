package repository

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// CreateUserKeycloak defines the interface to create a user in Keycloak.
type CreateUserKeycloak interface {
	// Create creates a new user in Keycloak.
	// It returns the ID of created user.
	Create(ctx context.Context, user *entity.User) (string, error)
}

// UserRegistrator is responsible to connect user with repositories.
type UserRegistrator struct {
	keycloak CreateUserKeycloak
	postgres service.RegisterUserRepository
}

// NewUserRegistrator creates an instance of UserRegistrator.
func NewUserRegistrator(keycloak CreateUserKeycloak, postgres service.RegisterUserRepository) *UserRegistrator {
	return &UserRegistrator{
		keycloak: keycloak,
		postgres: postgres,
	}
}

// Insert inserts user to Keycloak and Postgres.
func (ur *UserRegistrator) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}

	// TODO: keycloak and postgres must be atomic.
	id, err := ur.keycloak.Create(ctx, user)
	if err != nil {
		return err
	}
	user.KeycloakID = id
	if err := ur.postgres.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}
