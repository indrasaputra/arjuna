package service

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// RegisterUser defines the interface to register a user.
type RegisterUser interface {
	// Register registers a user and store it in the storage.
	// It returns the ID of the newly created user.
	//
	// It must check the uniqueness of the user.
	Register(ctx context.Context, user *entity.User) (string, error)
}
