package repository

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserRegistrator is responsible to connect user with repositories.
type UserRegistrator struct {
	keycloak service.RegisterUserRepository
	postgres service.RegisterUserRepository
}

// NewUserRegistrator creates an instance of UserRegistrator.
func NewUserRegistrator(keycloak, postgres service.RegisterUserRepository) *UserRegistrator {
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
	if err := ur.keycloak.Insert(ctx, user); err != nil {
		return err
	}
	if err := ur.postgres.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}
