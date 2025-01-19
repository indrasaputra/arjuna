package activity

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
)

// RegisterUserAuthConnection defines interface to register user to 3rd party.
type RegisterUserAuthConnection interface {
	// CreateAccount creates an account in 3rd party.
	CreateAccount(ctx context.Context, user *entity.User) error
}

// RegisterUserWalletConnection defines interface to register user to 3rd party.
type RegisterUserWalletConnection interface {
	// CreateWallet creates a wallet in 3rd party.
	CreateWallet(ctx context.Context, user *entity.User) error
}

// RegisterUserDatabase defines interface to register user to database.
type RegisterUserDatabase interface {
	// HardDelete hard-deletes the user. It must be called when a creation is failing and need a clean up or rollback.
	HardDelete(ctx context.Context, id uuid.UUID) error
}

// RegisterUserActivity is responsible to execute register user workflow.
type RegisterUserActivity struct {
	authConn   RegisterUserAuthConnection
	walletConn RegisterUserWalletConnection
	database   RegisterUserDatabase
}

// NewRegisterUserActivity creates an instance of RegisterUserActivity.
func NewRegisterUserActivity(ac RegisterUserAuthConnection, wc RegisterUserWalletConnection, db RegisterUserDatabase) *RegisterUserActivity {
	return &RegisterUserActivity{authConn: ac, walletConn: wc, database: db}
}

// CreateAccount creates a user in auth service.
func (r *RegisterUserActivity) CreateAccount(ctx context.Context, user *entity.User) error {
	err := r.authConn.CreateAccount(ctx, user)
	if errors.Is(err, entity.ErrAlreadyExists()) {
		return temporal.NewNonRetryableApplicationError(err.Error(), workflow.ErrNonRetryableUserExist, err)
	}
	return err
}

// CreateWallet creates user's wallet in wallet service.
func (r *RegisterUserActivity) CreateWallet(ctx context.Context, user *entity.User) error {
	err := r.walletConn.CreateWallet(ctx, user)
	if errors.Is(err, entity.ErrAlreadyExists()) {
		return temporal.NewNonRetryableApplicationError(err.Error(), workflow.ErrNonRetryableUserExist, err)
	}
	return err
}

// HardDeleteInUser hard-deletes user from database.
func (r *RegisterUserActivity) HardDeleteInUser(ctx context.Context, id uuid.UUID) error {
	err := r.database.HardDelete(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[RegisterUserActivity-HardDeleteInUser] fail hard delete in user", "error", err)
	}
	return err
}
