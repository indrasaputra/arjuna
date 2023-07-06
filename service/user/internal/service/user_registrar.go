package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

const (
	allowedUsernameCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
	usernameLength            = 10
)

var (
	regexNameCompiler *regexp.Regexp
)

func init() {
	regexNameCompiler = regexp.MustCompile(`^[a-zA-Z\s]+$`)
}

// RegisterUser defines the interface to register a user.
type RegisterUser interface {
	// Register registers a user and store it in the storage.
	// It returns the ID of the newly created user.
	// It must check the uniqueness of the user.
	Register(ctx context.Context, user *entity.User) (string, error)
}

// RegisterUserRepository defines interface to register user to repository.
type RegisterUserRepository interface {
	// InsertWithTx inserts user to repository using transaction.
	InsertWithTx(ctx context.Context, tx uow.Tx, user *entity.User) error
}

// RegisterUserOutboxRepository defines interface to register user outbox to repository.
type RegisterUserOutboxRepository interface {
	// InsertWithTx inserts user outbox to repository using transaction.
	InsertWithTx(ctx context.Context, tx uow.Tx, payload *entity.UserOutbox) error
}

// UserRegistrar is responsible for registering a new user.
type UserRegistrar struct {
	userRepo       RegisterUserRepository
	userOutboxRepo RegisterUserOutboxRepository
	unit           uow.UnitOfWork
}

// NewUserRegistrar creates an instance of UserRegistrar.
func NewUserRegistrar(ur RegisterUserRepository, uor RegisterUserOutboxRepository, u uow.UnitOfWork) *UserRegistrar {
	return &UserRegistrar{
		userRepo:       ur,
		userOutboxRepo: uor,
		unit:           u,
	}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRegistrar) Register(ctx context.Context, user *entity.User) (string, error) {
	if err := validateUser(user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] user is invalid: %v", err)
		return "", err
	}
	sanitizeUser(user)

	// username is mandatory for Keycloak, but not for this current business.
	// hence, generating a random username is fine.
	user.Username = generateUsername(usernameLength)

	if err := setUserID(user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail set user id: %v", err)
		return "", err
	}
	setUserAuditableProperties(user)

	err := ur.saveUserToRepository(ctx, user)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail save to repository: %v", err)
		return "", err
	}
	return user.ID, nil
}

func (ur *UserRegistrar) saveUserToRepository(ctx context.Context, user *entity.User) error {
	tx, err := ur.unit.Begin(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail init transaction: %v", err)
		return err
	}

	if err = ur.userRepo.InsertWithTx(ctx, tx, user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail insert user to repo: %v", err)
		_ = ur.unit.Finish(ctx, tx, err)
		return err
	}

	payload, err := createUserOutbox(user)
	if err != nil {
		_ = ur.unit.Finish(ctx, tx, err)
		return err
	}

	err = ur.userOutboxRepo.InsertWithTx(ctx, tx, payload)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail insert user outbox to repo: %v", err)
		_ = ur.unit.Finish(ctx, tx, err)
		return err
	}
	return ur.unit.Finish(ctx, tx, err)
}

func validateUser(user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}
	if !regexNameCompiler.MatchString(user.Name) {
		return entity.ErrInvalidName()
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return entity.ErrInvalidEmail()
	}
	return nil
}

func sanitizeUser(user *entity.User) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
}

func setUserID(user *entity.User) error {
	id, err := generateUniqueID()
	if err != nil {
		return entity.ErrInternal("fail to create user's ID")
	}
	user.ID = id
	return nil
}

func generateUniqueID() (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", entity.ErrInternal("fail to generate unique ID")
	}
	return id.String(), err
}

func setUserAuditableProperties(user *entity.User) {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.CreatedBy = user.ID
	user.UpdatedBy = user.ID
}

func generateUsername(n int) string {
	username := ""
	length := len(allowedUsernameCharacters)
	for i := 0; i < n; i++ {
		username += string(allowedUsernameCharacters[cryptoRandSecure(int64(length))])
	}
	return username
}

func cryptoRandSecure(max int64) int64 {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(max))
	return nBig.Int64()
}

func createUserOutbox(user *entity.User) (*entity.UserOutbox, error) {
	id, err := generateUniqueID()
	if err != nil {
		return nil, entity.ErrInternal("fail to generate user outbox id")
	}

	return &entity.UserOutbox{
		ID:      id,
		Status:  entity.UserOutboxStatusReady,
		Payload: user,
		Auditable: entity.Auditable{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedBy: user.UpdatedBy,
		},
	}, nil
}
