package handler

import (
	"context"
	"strings"

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

	clientID := request.GetCredential().GetClientId()
	email := request.GetCredential().GetEmail()
	password := request.GetCredential().GetPassword()

	token, err := a.auth.Login(ctx, clientID, email, password)
	if err != nil {
		app.Logger.Errorf(ctx, "[AuthHandler-Login] login fail: %v", err)
		return nil, err
	}
	return &apiv1.LoginResponse{Data: createTokenProto(token)}, nil
}

func validateLoginRequest(request *apiv1.LoginRequest) error {
	if request == nil || request.GetCredential() == nil {
		return entity.ErrEmptyField("request body")
	}
	if strings.TrimSpace(request.GetCredential().GetClientId()) == "" {
		return entity.ErrEmptyField("client id")
	}
	if strings.TrimSpace(request.GetCredential().GetEmail()) == "" {
		return entity.ErrEmptyField("email")
	}
	if strings.TrimSpace(request.GetCredential().GetPassword()) == "" {
		return entity.ErrEmptyField("password")
	}
	return nil
}

func createTokenProto(token *entity.Token) *apiv1.Token {
	return &apiv1.Token{
		AccessToken:           token.AccessToken,
		AccessTokenExpiresIn:  token.AccessTokenExpiresIn,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresIn: token.RefreshTokenExpiresIn,
		TokenType:             token.TokenType,
	}
}
