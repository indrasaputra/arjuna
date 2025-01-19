package auth

import (
	"context"
	"log/slog"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"

	enauth "github.com/indrasaputra/arjuna/service/auth/entity"
	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
	"github.com/indrasaputra/arjuna/service/user/entity"
)

// Auth is responsible to connect to auth service.
type Auth struct {
	client *sdkauth.Client
}

// NewAuth creates an instance of Auth.
func NewAuth(c *sdkauth.Client) *Auth {
	return &Auth{client: c}
}

// CreateAccount creates an account.
func (a *Auth) CreateAccount(ctx context.Context, user *entity.User) error {
	req := &enauth.Account{UserID: user.ID, Email: user.Email, Password: user.Password}
	err := a.client.Register(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[Auth-CreateAccount] fail call register", "error", err)
	}
	if status.Code(err) == codes.AlreadyExists {
		return entity.ErrAlreadyExists()
	}
	return err
}
