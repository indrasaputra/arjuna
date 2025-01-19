package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	limitGetAllReady = uint(1)
)

// RelayRegisterUser defines the interface to relay the user registration.
type RelayRegisterUser interface {
	// RelayRegister relays user registration.
	RelayRegister(ctx context.Context, user *entity.UserOutbox) (string, error)
}

// RelayRegisterUserOrchestration defines the interface to orchestrate user registration.
type RelayRegisterUserOrchestration interface {
	// RegisterUser registers the user into the repository or 3rd party needed.
	// It also validates if the user's email is unique.
	// It returns the ID of the created user.
	RegisterUser(ctx context.Context, input *entity.RegisterUserInput) (*entity.RegisterUserOutput, error)
}

// RelayRegisterUserOutboxRepository defines interface to register user outbox to repository.
type RelayRegisterUserOutboxRepository interface {
	// GetAllReady gets all ready records.
	GetAllReady(ctx context.Context, limit uint) ([]*entity.UserOutbox, error)
	// SetProcessed sets record as processed.
	SetProcessed(ctx context.Context, id uuid.UUID) error
	// SetDelivered sets record as delivered.
	SetDelivered(ctx context.Context, id uuid.UUID) error
	// SetFailed sets record as failed.
	SetFailed(ctx context.Context, id uuid.UUID) error
}

// UserRelayRegistrar is responsible for registering a new user.
type UserRelayRegistrar struct {
	userOutboxRepo RelayRegisterUserOutboxRepository
	orchestrator   RelayRegisterUserOrchestration
	txManager      uow.TxManager
}

// NewUserRelayRegistrar creates an instance of UserRelayRegistrar.
func NewUserRelayRegistrar(r RelayRegisterUserOutboxRepository, o RelayRegisterUserOrchestration, t uow.TxManager) *UserRelayRegistrar {
	return &UserRelayRegistrar{
		userOutboxRepo: r,
		orchestrator:   o,
		txManager:      t,
	}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRelayRegistrar) Register(ctx context.Context) error {
	err := ur.txManager.Do(ctx, func(ctx context.Context) error {
		records, err := ur.userOutboxRepo.GetAllReady(ctx, limitGetAllReady)
		if err != nil {
			slog.ErrorContext(ctx, "[UserRelayRegistrar-Register] fail get all ready", "error", err)
			return err
		}

		for _, rc := range records {
			if err := ur.setRecordAsProcessed(ctx, rc); err != nil {
				slog.ErrorContext(ctx, "[UserRelayRegistrar-Register] fail set record as processed", "record", rc, "error", err)
				continue
			}
			if err := ur.enqueueRecordToOrchestrator(ctx, rc); err != nil {
				slog.ErrorContext(ctx, "[UserRelayRegistrar-Register] fail enqueue record", "record", rc, "error", err)
			}
		}
		return nil
	})
	return err
}

func (ur *UserRelayRegistrar) enqueueRecordToOrchestrator(ctx context.Context, record *entity.UserOutbox) error {
	err := ur.startRegisterUserWorkflow(ctx, record)
	if err != nil {
		slog.ErrorContext(ctx, "[UserRelayRegistrar-Register] fail start workflow", "record", record, "error", err)
		return ur.setRecordAsFailed(ctx, record)
	}
	return ur.setRecordAsDelivered(ctx, record)
}

func (ur *UserRelayRegistrar) startRegisterUserWorkflow(ctx context.Context, record *entity.UserOutbox) error {
	input := &entity.RegisterUserInput{User: record.Payload}
	_, err := ur.orchestrator.RegisterUser(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "[UserRelayRegistrar-Register] orchestration fail", "error", err)
		return err
	}
	return nil
}

func (ur *UserRelayRegistrar) setRecordAsProcessed(ctx context.Context, record *entity.UserOutbox) error {
	return ur.userOutboxRepo.SetProcessed(ctx, record.ID)
}

func (ur *UserRelayRegistrar) setRecordAsDelivered(ctx context.Context, record *entity.UserOutbox) error {
	return ur.userOutboxRepo.SetDelivered(ctx, record.ID)
}

func (ur *UserRelayRegistrar) setRecordAsFailed(ctx context.Context, record *entity.UserOutbox) error {
	return ur.userOutboxRepo.SetFailed(ctx, record.ID)
}
