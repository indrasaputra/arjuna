package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
)

// Auth handles HTTP/2 gRPC request for auth.
type Auth struct {
	apiv1.UnimplementedAuthServiceServer
	auth service.Authentication
}

// NewAuth creates an instance of Auth.
func NewAuth(auth service.Authentication) *Auth {
	return &Auth{auth: auth}
}

// Login handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (a *Auth) Login(ctx context.Context, request *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	if err := validateLoginRequest(request); err != nil {
		app.Logger.Errorf(ctx, "[AuthHandler-Login] request invalid: %v", err)
		return nil, err
	}

	email := request.GetCredential().GetEmail()
	password := request.GetCredential().GetPassword()

	token, err := a.auth.Login(ctx, email, password)
	if err != nil {
		app.Logger.Errorf(ctx, "[AuthHandler-Login] login fail: %v", err)
		return nil, err
	}
	return &apiv1.LoginResponse{Data: createTokenProto(token)}, nil
}

// RegisterAccount handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (a *Auth) RegisterAccount(ctx context.Context, request *apiv1.RegisterAccountRequest) (*apiv1.RegisterAccountResponse, error) {
	if err := validateRegisterAccountRequest(request); err != nil {
		app.Logger.Errorf(ctx, "[AuthHandler-Register] request invalid: %v", err)
		return nil, err
	}

	err := a.auth.Register(ctx, createAccountFromRegisterAccountRequest(request))
	if err != nil {
		app.Logger.Errorf(ctx, "[AuthHandler-Register] login fail: %v", err)
		return nil, err
	}
	return &apiv1.RegisterAccountResponse{}, nil
}

func validateLoginRequest(request *apiv1.LoginRequest) error {
	if request == nil || request.GetCredential() == nil {
		return entity.ErrEmptyField("request body")
	}
	if strings.TrimSpace(request.GetCredential().GetEmail()) == "" {
		return entity.ErrEmptyField("email")
	}
	if strings.TrimSpace(request.GetCredential().GetPassword()) == "" {
		return entity.ErrEmptyField("password")
	}
	return nil
}

func validateRegisterAccountRequest(request *apiv1.RegisterAccountRequest) error {
	if request == nil || request.GetAccount() == nil {
		return entity.ErrEmptyField("request body")
	}
	if strings.TrimSpace(request.GetAccount().GetUserId()) == "" {
		return entity.ErrEmptyField("user id")
	}
	if strings.TrimSpace(request.GetAccount().GetEmail()) == "" {
		return entity.ErrEmptyField("email")
	}
	if strings.TrimSpace(request.GetAccount().GetPassword()) == "" {
		return entity.ErrEmptyField("password")
	}
	return nil
}

func createAccountFromRegisterAccountRequest(request *apiv1.RegisterAccountRequest) *entity.Account {
	return &entity.Account{
		UserID:   uuid.MustParse(request.GetAccount().GetUserId()),
		Email:    request.GetAccount().GetEmail(),
		Password: request.GetAccount().GetPassword(),
	}
}

func createTokenProto(token *entity.Token) *apiv1.Token {
	return &apiv1.Token{
		AccessToken:           token.AccessToken,
		AccessTokenExpiresIn:  token.AccessTokenExpiresIn,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresIn: token.RefreshTokenExpiresIn,
	}
}
