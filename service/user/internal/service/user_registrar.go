package service

import (
	"context"
	"log/slog"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
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
	Register(ctx context.Context, user *entity.User) (uuid.UUID, error)
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

// UserRegistrar is responsible for registering a new user.
type UserRegistrar struct {
	txManager      uow.TxManager
	userRepo       RegisterUserRepository
	userOutboxRepo RegisterUserOutboxRepository
}

// NewUserRegistrar creates an instance of UserRegistrar.
func NewUserRegistrar(txm uow.TxManager, ur RegisterUserRepository, uor RegisterUserOutboxRepository) *UserRegistrar {
	return &UserRegistrar{
		txManager:      txm,
		userRepo:       ur,
		userOutboxRepo: uor,
	}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRegistrar) Register(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	sanitizeUser(user)
	if err := validateUser(user); err != nil {
		slog.ErrorContext(ctx, "[UserRegistrar-Register] fail validate user", "error", err)
		return uuid.Nil, err
	}

	setUserID(user)
	setUserAuditableProperties(user)

	err := ur.saveUserToRepository(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "[UserRegistrar-Register] fail save to repository", "error", err)
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (ur *UserRegistrar) saveUserToRepository(ctx context.Context, user *entity.User) error {
	err := ur.txManager.Do(ctx, func(ctx context.Context) error {
		if err := ur.userRepo.Insert(ctx, user); err != nil {
			slog.ErrorContext(ctx, "[UserRegistrar-saveUserToRepository] fail insert user to repo", "error", err)
			return err
		}
		payload := createUserOutbox(user)
		err := ur.userOutboxRepo.Insert(ctx, payload)
		if err != nil {
			slog.ErrorContext(ctx, "[UserRegistrar-saveUserToRepository] fail insert user outbox to repo", "error", err)
		}
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "[UserRegistrar-saveUserToRepository] transaction fail", "error", err)
	}
	return err
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
	user.CreatedBy = user.ID
	user.UpdatedBy = user.ID
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
