package activity

import (
	"context"
	"errors"

	"go.temporal.io/sdk/temporal"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
)

// RegisterUserVendor defines interface to register user to vendor.
type RegisterUserVendor interface {
	// Create creates a user in vendor. It must returns the user's ID in vendor's side.
	Create(ctx context.Context, user *entity.User) (string, error)
	// HardDelete hard-deletes the user. It must be called when a creation is failing and need a clean up or rollback.
	HardDelete(ctx context.Context, id string) error
}

// RegisterUserDatabase defines interface to register user to database.
type RegisterUserDatabase interface {
	// Insert inserts a user into a database.
	Insert(ctx context.Context, user *entity.User) error
}

// RegisterUserActivity is responsible to execute register user workflow.
type RegisterUserActivity struct {
	vendor   RegisterUserVendor
	database RegisterUserDatabase
}

// NewRegisterUserActivity creates an instance of RegisterUserActivity.
func NewRegisterUserActivity(pvd RegisterUserVendor, db RegisterUserDatabase) *RegisterUserActivity {
	return &RegisterUserActivity{vendor: pvd, database: db}
}

// CreateInKeycloak creates a user in Keycloak (vendor).
func (r *RegisterUserActivity) CreateInKeycloak(ctx context.Context, user *entity.User) (string, error) {
	id, err := r.vendor.Create(ctx, user)
	if errors.Is(err, entity.ErrAlreadyExists()) {
		return "", temporal.NewNonRetryableApplicationError(err.Error(), workflow.ErrNonRetryableUserExist, err)
	}
	return id, err
}

// HardDeleteFromKeycloak hard-deletes user from Keycloak.
func (r *RegisterUserActivity) HardDeleteFromKeycloak(ctx context.Context, id string) error {
	return r.vendor.HardDelete(ctx, id)
}

// InsertToDatabase inserts user to database.
func (r *RegisterUserActivity) InsertToDatabase(ctx context.Context, user *entity.User) error {
	err := r.database.Insert(ctx, user)
	if errors.Is(err, entity.ErrAlreadyExists()) {
		return temporal.NewNonRetryableApplicationError(err.Error(), workflow.ErrNonRetryableUserExist, err)
	}
	return err
}
