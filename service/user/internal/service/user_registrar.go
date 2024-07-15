package service

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"

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
	Register(ctx context.Context, user *entity.User, key string) (string, error)
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

// IdempotencyKeyRepository defines  interface for idempotency check flow and repository.
type IdempotencyKeyRepository interface {
	// Exists check if given key exists in repository.
	Exists(ctx context.Context, key string) (bool, error)
}

// UserRegistrar is responsible for registering a new user.
type UserRegistrar struct {
	userRepo       RegisterUserRepository
	userOutboxRepo RegisterUserOutboxRepository
	unit           uow.UnitOfWork
	keyRepo        IdempotencyKeyRepository
}

// NewUserRegistrar creates an instance of UserRegistrar.
func NewUserRegistrar(ur RegisterUserRepository, uor RegisterUserOutboxRepository, u uow.UnitOfWork, k IdempotencyKeyRepository) *UserRegistrar {
	return &UserRegistrar{
		userRepo:       ur,
		userOutboxRepo: uor,
		unit:           u,
		keyRepo:        k,
	}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRegistrar) Register(ctx context.Context, user *entity.User, key string) (string, error) {
	if err := ur.validateIdempotencyKey(ctx, key); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail check idempotency key: %s - %v", key, err)
		return "", err
	}

	sanitizeUser(user)
	if err := validateUser(user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] user is invalid: %v", err)
		return "", err
	}

	if err := setUserID(ctx, user); err != nil {
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
		return entity.ErrInternal("fail to begin transaction")
	}

	if err = ur.userRepo.InsertWithTx(ctx, tx, user); err != nil {
		app.Logger.Errorf(ctx, "[UserRegistrar-saveUserToRepository] fail insert user to repo: %v", err)
		_ = ur.unit.Finish(ctx, tx, err)
		return err
	}

	payload, err := createUserOutbox(ctx, user)
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

func setUserID(ctx context.Context, user *entity.User) error {
	id, err := generateUniqueID(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[setUserID] fail generate unique id: %v", err)
		return entity.ErrInternal("fail to create user's ID")
	}
	user.ID = id
	return nil
}

func generateUniqueID(ctx context.Context) (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		app.Logger.Errorf(ctx, "[setUserID] fail generate ksuid: %v", err)
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

func createUserOutbox(ctx context.Context, user *entity.User) (*entity.UserOutbox, error) {
	id, err := generateUniqueID(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[createUserOutbox] fail generate unique id: %v", err)
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
