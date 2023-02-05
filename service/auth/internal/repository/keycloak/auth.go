package keycloak

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
)

// Config defines Keycloak config.
type Config struct {
	// Client is required.
	Client kcsdk.Keycloak
	// Realm is required.
	Realm string
}

// Validate validates all fields in config.
// For example, all required fields must be set.
func (c *Config) Validate() error {
	if c.Client == nil {
		return errors.New("client must be set")
	}
	if strings.TrimSpace(c.Realm) == "" {
		return errors.New("realm must be set")
	}
	return nil
}

// Auth is responsible to connect user entity with Keycloak.
type Auth struct {
	config *Config
}

// NewAuth creates an instance of Auth.
func NewAuth(config *Config) (*Auth, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &Auth{config: config}, nil
}

// Login logs in a user and returns token.
func (u *Auth) Login(ctx context.Context, clientID, email, password string) (*entity.Token, error) {
	jwt, err := u.config.Client.LoginUser(ctx, u.config.Realm, clientID, email, password)
	if err != nil {
		app.Logger.Errorf(ctx, "[AuthRepo-Login] login fail: %v", err)
		return nil, decideError(err)
	}
	return createToken(jwt), nil
}

func createToken(jwt *kcsdk.JWT) *entity.Token {
	return &entity.Token{
		AccessToken:           jwt.AccessToken,
		AccessTokenExpiresIn:  uint32(jwt.ExpiresIn),
		RefreshToken:          jwt.RefreshToken,
		RefreshTokenExpiresIn: uint32(jwt.RefreshExpiresIn),
		TokenType:             jwt.TokenType,
	}
}

func decideError(err error) error {
	kcerr, ok := err.(*kcsdk.Error)
	if !ok {
		log.Printf("[Auth-Login] unknown error from keycloak: %v\n", err)
		return entity.ErrInternal(err.Error())
	}

	switch kcerr.Code {
	case http.StatusUnauthorized:
		return entity.ErrUnauthorized()
	case http.StatusBadRequest:
		return entity.ErrInvalidArgument(kcerr.Error())
	default:
		return entity.ErrInternal(err.Error())
	}
}
