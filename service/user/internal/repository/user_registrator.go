package repository

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/internal/service"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// // RegisterUserRepository defines the interface to save a user into the repository.
// type RegisterUserRepository interface {
// 	// Insert inserts the user into the repository.
// 	// It also validates if the user's email is unique.
// 	Insert(ctx context.Context, user *entity.User) error
// }

type 

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

func (ur *UserRegistrator) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}


}
