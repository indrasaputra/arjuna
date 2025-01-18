package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
)

const (
	tokenIssuer = "auth-service"
	timeMinute  = 60
)

// Authentication defines the interface to authenticate.
type Authentication interface {
	// Login logs in a user using email and password.
	Login(ctx context.Context, email, password string) (*entity.Token, error)
	// Register registers an account.
	Register(ctx context.Context, account *entity.Account) error
}

// AuthRepository defines the interface to authenticate.
type AuthRepository interface {
	// GetByEmail gets an account by email.
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	// Insert inserts an account.
	Insert(ctx context.Context, account *entity.Account) error
}

// Auth is responsible for authentication.
type Auth struct {
	repo            AuthRepository
	signingKey      []byte
	tokenExpiration int
}

// NewAuth creates an instance of Auth.
func NewAuth(repo AuthRepository, key []byte, exp int) *Auth {
	return &Auth{repo: repo, tokenExpiration: exp, signingKey: key}
}

// Login logs in a user using email and password.
// As of now, refresh token is not implemented and it only returns access token.
func (a *Auth) Login(ctx context.Context, email, password string) (*entity.Token, error) {
	if err := validateLoginParams(email, password); err != nil {
		app.Logger.Errorf(ctx, "[Auth-Login] param invalid: %v", err)
		return nil, err
	}

	account, err := a.repo.GetByEmail(ctx, email)
	if status.Code(err) == codes.NotFound {
		return nil, entity.ErrInvalidCredential()
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return nil, entity.ErrInvalidCredential()
	}
	return createAccessToken(account, a.signingKey, a.tokenExpiration)
}

// Register registers an account.
func (a *Auth) Register(ctx context.Context, account *entity.Account) error {
	sanitizeAccount(account)
	if err := validateAccount(account); err != nil {
		app.Logger.Errorf(ctx, "[Auth-Register] account is invalid: %v", err)
		return err
	}

	setAccountID(account)
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
	if account == nil || account.UserID == uuid.Nil {
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
	account.Email = strings.TrimSpace(account.Email)
	account.Password = strings.TrimSpace(account.Password)
}

func validateLoginParams(email, password string) error {
	if strings.TrimSpace(email) == "" {
		return entity.ErrEmptyField("email")
	}
	if strings.TrimSpace(password) == "" {
		return entity.ErrEmptyField("password")
	}
	return nil
}

func setAccountID(account *entity.Account) {
	account.ID = generateUniqueID()
}

func generateUniqueID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
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

func createAccessToken(account *entity.Account, key []byte, exp int) (*entity.Token, error) {
	claims := entity.Claims{
		AccountID: account.ID,
		UserID:    account.UserID,
		Email:     account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exp) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tokenIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString(key)
	if err != nil {
		return nil, entity.ErrInternal("fail signing token")
	}
	return &entity.Token{AccessToken: res, AccessTokenExpiresIn: uint32(exp * timeMinute)}, nil
}
