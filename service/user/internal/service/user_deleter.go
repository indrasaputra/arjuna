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
	// GetByID gets a user by ID.
	GetByID(ctx context.Context, id string) (*entity.User, error)
	// HardDelete hard-deletes a single user from the repository.
	// If the user can't be found, it doesn't return error.
	HardDelete(ctx context.Context, user *entity.User) error
}

// UserDeleter is responsible for deleting a user.
type UserDeleter struct {
	repo DeleteUserRepository
}

// NewUserDeleter creates an instance of UserDeleter.
func NewUserDeleter(repo DeleteUserRepository) *UserDeleter {
	return &UserDeleter{
		repo: repo,
	}
}

// HardDelete hard-deletes a user in system.
func (td *UserDeleter) HardDelete(ctx context.Context, id string) error {
	user, err := td.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return td.repo.HardDelete(ctx, user)
}
