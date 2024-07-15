package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
)

// Authentication defines the interface to authenticate.
type Authentication interface {
	// Login logs in a user using email and password.
	Login(ctx context.Context, clientID, email, password string) (*entity.Token, error)
	// Register registers an account.
	Register(ctx context.Context, account *entity.Account) error
}

// AuthRepository defines the interface to authenticate.
type AuthRepository interface {
	// Login logs in a user.
	Login(ctx context.Context, clientID, email, password string) (*entity.Token, error)
	// Insert inserts an account.
	Insert(ctx context.Context, account *entity.Account) error
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
	if err := validateLoginParams(clientID, email, password); err != nil {
		app.Logger.Errorf(ctx, "[Auth-Login] param invalid: %v", err)
		return nil, err
	}
	return a.repo.Login(ctx, clientID, email, password)
}

// Register registers an account.
func (a *Auth) Register(ctx context.Context, account *entity.Account) error {
	sanitizeAccount(account)
	if err := validateAccount(account); err != nil {
		app.Logger.Errorf(ctx, "[Auth-Register] account is invalid: %v", err)
		return err
	}

	if err := setAccountID(ctx, account); err != nil {
		app.Logger.Errorf(ctx, "[Auth-Register] fail set account id: %v", err)
		return err
	}
	setAccountAuditableProperties(account)

	hash, err := encryptPassword(ctx, account.Password)
	if err != nil {
		return err
	}
	account.Password = hash

	err = a.repo.Insert(ctx, account)
	if err != nil {
		app.Logger.Errorf(ctx, "[Auth-Register] fail save to repository: %v", err)
		return err
	}
	return nil
}

func validateAccount(account *entity.Account) error {
	if account == nil || account.UserID == "" {
		return entity.ErrEmptyAccount()
	}
	if _, err := mail.ParseAddress(account.Email); err != nil {
		return entity.ErrInvalidEmail()
	}
	if account.Password == "" {
		return entity.ErrInvalidPassword()
	}
	return nil
}

func sanitizeAccount(account *entity.Account) {
	if account == nil {
		return
	}
	account.UserID = strings.TrimSpace(account.UserID)
	account.Email = strings.TrimSpace(account.Email)
	account.Password = strings.TrimSpace(account.Password)
}

func validateLoginParams(clientID, email, password string) error {
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

func setAccountID(ctx context.Context, account *entity.Account) error {
	id, err := generateUniqueID(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[setAccountID] fail generate unique id: %v", err)
		return entity.ErrInternal("fail to create account's ID")
	}
	account.ID = id
	return nil
}

func generateUniqueID(ctx context.Context) (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		app.Logger.Errorf(ctx, "[setAccountID] fail generate ksuid: %v", err)
		return "", entity.ErrInternal("fail to generate unique ID")
	}
	return id.String(), err
}

func encryptPassword(ctx context.Context, password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.Logger.Errorf(ctx, "[encryptPassword] fail generate from password: %v", err)
		return "", entity.ErrInternal("fail to hash the password")
	}
	return string(res), nil
}

func setAccountAuditableProperties(account *entity.Account) {
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = time.Now().UTC()
	account.CreatedBy = account.ID
	account.UpdatedBy = account.ID
}
