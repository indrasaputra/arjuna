package activity

import (
	"context"
	"errors"

	"go.temporal.io/sdk/temporal"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
)

// RegisterUserConnection defines interface to register user to 3rd party.
type RegisterUserConnection interface {
	// CreateAccount creates an account in 3rd party.
	CreateAccount(ctx context.Context, user *entity.User) error
}

// RegisterUserDatabase defines interface to register user to database.
type RegisterUserDatabase interface {
	// HardDelete hard-deletes the user. It must be called when a creation is failing and need a clean up or rollback.
	HardDelete(ctx context.Context, id string) error
}

// RegisterUserActivity is responsible to execute register user workflow.
type RegisterUserActivity struct {
	conn     RegisterUserConnection
	database RegisterUserDatabase
}

// NewRegisterUserActivity creates an instance of RegisterUserActivity.
func NewRegisterUserActivity(conn RegisterUserConnection, db RegisterUserDatabase) *RegisterUserActivity {
	return &RegisterUserActivity{conn: conn, database: db}
}

// CreateInAuth creates a user in auth service.
func (r *RegisterUserActivity) CreateInAuth(ctx context.Context, user *entity.User) error {
	err := r.conn.CreateAccount(ctx, user)
	if errors.Is(err, entity.ErrAlreadyExists()) {
		return temporal.NewNonRetryableApplicationError(err.Error(), workflow.ErrNonRetryableUserExist, err)
	}
	return err
}

// HardDeleteInUser hard-deletes user from database.
func (r *RegisterUserActivity) HardDeleteInUser(ctx context.Context, id string) error {
	err := r.database.HardDelete(ctx, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[RegisterUserActivity-HardDeleteInUser] fail hard delete in user: %v", err)
	}
	return err
}
