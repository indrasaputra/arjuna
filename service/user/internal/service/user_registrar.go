package service

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
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
	Register(ctx context.Context, user *entity.User, key string) (uuid.UUID, error)
}

// RegisterUserRepository defines interface to register user to repository.
type RegisterUserRepository interface {
	// InsertWithTx inserts user to repository using transaction.
	// InsertWithTx(ctx context.Context, tx uow.Tx, user *entity.User) error
	Insert(ctx context.Context, user *entity.User) error
}

// RegisterUserOutboxRepository defines interface to register user outbox to repository.
type RegisterUserOutboxRepository interface {
	// InsertWithTx inserts user outbox to repository using transaction.
	// InsertWithTx(ctx context.Context, tx uow.Tx, payload *entity.UserOutbox) error
	Insert(ctx context.Context, payload *entity.UserOutbox) error
}

// IdempotencyKeyRepository defines  interface for idempotency check flow and repository.
type IdempotencyKeyRepository interface {
	// Exists check if given key exists in repository.
	Exists(ctx context.Context, key string) (bool, error)
}

// UserRegistrar is responsible for registering a new user.
type UserRegistrar struct {
	txManager      uow.TxManager
	userRepo       RegisterUserRepository
	userOutboxRepo RegisterUserOutboxRepository
	keyRepo        IdempotencyKeyRepository
}

// NewUserRegistrar creates an instance of UserRegistrar.
func NewUserRegistrar(txm uow.TxManager, ur RegisterUserRepository, uor RegisterUserOutboxRepository, k IdempotencyKeyRepository) *UserRegistrar {
	return &UserRegistrar{
		txManager:      txm,
		userRepo:       ur,
		userOutboxRepo: uor,
		keyRepo:        k,
	}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRegistrar) Register(ctx context.Context, user *entity.User, key string) (uuid.UUID, error) {
	if err := ur.validateIdempotencyKey(ctx, key); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail check idempotency key: %s - %v", key, err)
		return uuid.Nil, err
	}

	sanitizeUser(user)
	if err := validateUser(user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] user is invalid: %v", err)
		return uuid.Nil, err
	}

	setUserID(user)
	setUserAuditableProperties(user)

	err := ur.saveUserToRepository(ctx, user)
	if err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail save to repository: %v", err)
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (ur *UserRegistrar) saveUserToRepository(ctx context.Context, user *entity.User) error {
	err := ur.txManager.Do(ctx, func(ctx context.Context) error {
		if err := ur.userRepo.Insert(ctx, user); err != nil {
			app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail insert user to repo: %v", err)
			return err
		}
		payload := createUserOutbox(user)
		err := ur.userOutboxRepo.Insert(ctx, payload)
		if err != nil {
			app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail insert user outbox to repo: %v", err)
		}
		return err
	})
	if err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] transaction fail: %v", err)
	}
	return err
}

func (ur *UserRegistrar) validateIdempotencyKey(ctx context.Context, key string) error {
	res, err := ur.keyRepo.Exists(ctx, key)
	if err != nil {
		return err
	}
	if res {
		return entity.ErrAlreadyExists()
	}
	return nil
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
	if user == nil {
		return
	}
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)
}

func setUserID(user *entity.User) {
	user.ID = generateUniqueID()
}

func generateUniqueID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func setUserAuditableProperties(user *entity.User) {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.CreatedBy = user.ID.String()
	user.UpdatedBy = user.ID.String()
}

func createUserOutbox(user *entity.User) *entity.UserOutbox {
	return &entity.UserOutbox{
		ID:      generateUniqueID(),
		Status:  entity.UserOutboxStatusReady,
		Payload: user,
		Auditable: entity.Auditable{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedBy: user.UpdatedBy,
		},
	}
}
