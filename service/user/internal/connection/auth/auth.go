package auth

import (
	"context"

	"github.com/gogo/status"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	enauth "github.com/indrasaputra/arjuna/service/auth/entity"
	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
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
	req := &enauth.Account{UserID: uuid.MustParse(user.ID), Email: user.Email, Password: user.Password}
	err := a.client.Register(ctx, req)
	if err != nil {
		app.Logger.Errorf(ctx, "[Auth-CreateAccount] fail call register: %v", err)
	}
	if status.Code(err) == codes.AlreadyExists {
		return entity.ErrAlreadyExists()
	}
	return err
}
