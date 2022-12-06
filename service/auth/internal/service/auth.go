package service

import (
	"context"
	"strings"

	"github.com/indrasaputra/arjuna/service/auth/entity"
)

// Authentication defines the interface to authenticate.
type Authentication interface {
	// Login logs in a user using email and password.
	Login(ctx context.Context, clientID, email, password string) (*entity.Token, error)
}

// AuthRepository defines the interface to authenticate.
type AuthRepository interface {
	// Login logs in a user.
	Login(ctx context.Context, clientID, email, password string) (*entity.Token, error)
}

// Auth is responsible for authentication.
type Auth struct {
	repo AuthRepository
}

// NewAuth creates an instance of Auth.
func NewAuth(repo AuthRepository) *Auth {
	return &Auth{repo: repo}
}

// Login logs in a user using email and password.
func (a *Auth) Login(ctx context.Context, clientID, email, password string) (*entity.Token, error) {
	if err := validateParams(clientID, email, password); err != nil {
		return nil, err
	}
	return a.repo.Login(ctx, clientID, email, password)
}

func validateParams(clientID, email, password string) error {
	if strings.TrimSpace(clientID) == "" {
		return entity.ErrEmptyField("clientId")
	}
	if strings.TrimSpace(email) == "" {
		return entity.ErrEmptyField("email")
	}
	if strings.TrimSpace(password) == "" {
		return entity.ErrEmptyField("password")
	}
	return nil
}
