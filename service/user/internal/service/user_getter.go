package service

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// GetUser defines the interface to get user.
type GetUser interface {
	// GetAll gets all users.
	GetAll(ctx context.Context) ([]*entity.User, error)
}

// GetUserRepository defines the interface to get user from the repository.
type GetUserRepository interface {
	// GetAll gets all users available in repository.
	// If there isn't any user in repository, it returns empty list of user and nil error.
	GetAll(ctx context.Context) ([]*entity.User, error)
}

// UserGetter is responsible for getting user.
type UserGetter struct {
	repo GetUserRepository
}

// NewUserGetter creates an instance of UserGetter.
func NewUserGetter(repo GetUserRepository) *UserGetter {
	return &UserGetter{repo: repo}
}

// GetAll gets all users available in repository.
// If there isn't any user in repository, it returns empty list of user and nil error.
func (ug *UserGetter) GetAll(ctx context.Context) ([]*entity.User, error) {
	return ug.repo.GetAll(ctx)
}
