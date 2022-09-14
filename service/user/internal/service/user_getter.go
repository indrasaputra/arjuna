package service

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	defaultLimit = 10
	// DefaultGetAllUsersLimit is deliberately set to 10 because it is enough for simple application.
	DefaultGetAllUsersLimit = uint(defaultLimit)
)

// GetUser defines the interface to get user.
type GetUser interface {
	// GetAll gets all users.
	GetAll(ctx context.Context, limit uint) ([]*entity.User, error)
}

// GetUserRepository defines the interface to get user from the repository.
type GetUserRepository interface {
	// GetAll gets all users available in repository.
	// If there isn't any user in repository, it returns empty list of user and nil error.
	GetAll(ctx context.Context, limit uint) ([]*entity.User, error)
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
func (ug *UserGetter) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	if limit == 0 || limit > DefaultGetAllUsersLimit {
		limit = DefaultGetAllUsersLimit
	}
	return ug.repo.GetAll(ctx, limit)
}
