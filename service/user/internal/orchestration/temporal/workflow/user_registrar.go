package workflow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

const (
	// TaskQueueRegisterUser represents user registration.
	TaskQueueRegisterUser = "register-user"

	// ActivityTimeoutDefault sets to 2 seconds.
	ActivityTimeoutDefault = 2 * time.Second
	// ActivityKeycloakCreate is derived from struct name + method name. See activity registration in worker.
	ActivityKeycloakCreate = "RegisterUserActivityCreateInKeycloak"
	// ActivityKeycloakHardDelete is derived from struct name + method name. See activity registration in worker.
	ActivityKeycloakHardDelete = "RegisterUserActivityHardDeleteFromKeycloak"
	// ActivityPostgresInsert is derived from struct name + method name. See activity registration in worker.
	ActivityPostgresInsert = "RegisterUserActivityInsertToDatabase"
	// ActivityRetryBackoffCoefficient sets to 2.
	ActivityRetryBackoffCoefficient = 2
	// ActivityRetryMaximumAttempts sets to 3.
	ActivityRetryMaximumAttempts = 3
	// ActivityRetryInitialInterval sets to 1 second.
	ActivityRetryInitialInterval = 1 * time.Second

	// WorkflowTimeoutDefault sets ActivityRetryMaximumAttempts * ActivityTimeoutDefault.
	WorkflowTimeoutDefault = ActivityRetryMaximumAttempts * ActivityTimeoutDefault
	// WorkflowNameRegisterUser is derived from the process itself.
	WorkflowNameRegisterUser = "register-user"
	// WorkflowRetryMaximumAttempts sets to 1.
	WorkflowRetryMaximumAttempts = 1

	// ErrNonRetryableUserExist occurs when user already exists in system.
	ErrNonRetryableUserExist = "non-retryable-user-exist"
)

// RegisterUserWorkflow is responsible to execute register user workflow.
type RegisterUserWorkflow struct {
	client client.Client
}

// NewRegisterUserWorkflow creates an instance of RegisterUserWorkflow.
func NewRegisterUserWorkflow(client client.Client) *RegisterUserWorkflow {
	return &RegisterUserWorkflow{client: client}
}

// RegisterUser runs the register users workflow.
func (r *RegisterUserWorkflow) RegisterUser(ctx context.Context, input *service.RegisterUserInput) (*service.RegisterUserOutput, error) {
	opts := client.StartWorkflowOptions{
		ID:                 fmt.Sprintf("%s-%s", WorkflowNameRegisterUser, input.User.ID),
		TaskQueue:          TaskQueueRegisterUser,
		WorkflowRunTimeout: WorkflowTimeoutDefault,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: WorkflowRetryMaximumAttempts,
			NonRetryableErrorTypes: []string{
				ErrNonRetryableUserExist,
			},
		},
	}
	wr, err := r.client.ExecuteWorkflow(ctx, opts, RegisterUser, input)
	if err != nil {
		app.Logger.Errorf(ctx, "[RegisterUserWorkflow-RegisterUser] workflow failure: %v", err)
		return nil, entity.ErrInternal("Something went wrong within our server. Please, try again")
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	var output *service.RegisterUserOutput
	err = wr.Get(ctx, &output)
	if err != nil {
		var appErr *temporal.ApplicationError
		if errors.As(err, &appErr) && appErr.Type() == ErrNonRetryableUserExist {
			return nil, entity.ErrAlreadyExists()
		}
		app.Logger.Errorf(ctx, "[RegisterUserWorkflow-RegisterUser] error get workflow result: %v", err)
		return nil, entity.ErrInternal("Something went wrong within our server. Please, try again")
	}
	return output, nil
}

// RegisterUser runs the user registration workflow.
func RegisterUser(ctx workflow.Context, input *service.RegisterUserInput) (*service.RegisterUserOutput, error) {
	if err := validateRegisterUserInput(input); err != nil {
		return nil, err
	}

	var id string
	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err := workflow.ExecuteActivity(ctx, ActivityKeycloakCreate, input.User).Get(ctx, &id)
	if err != nil {
		return nil, err
	}
	input.User.KeycloakID = id

	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err = workflow.ExecuteActivity(ctx, ActivityPostgresInsert, input.User).Get(ctx, nil)
	if err != nil {
		ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
		_ = workflow.ExecuteActivity(ctx, ActivityKeycloakHardDelete, id).Get(ctx, nil)
		return nil, entity.ErrInternal("Something went wrong within our server. Please try again")
	}
	return &service.RegisterUserOutput{UserID: input.User.ID}, nil
}

func createContextWithActivityOptions(tempoCtx workflow.Context, timeout time.Duration, queue string) workflow.Context {
	opts := createActivityOptions(timeout, queue)
	return workflow.WithActivityOptions(tempoCtx, opts)
}

func createActivityOptions(timeout time.Duration, queue string) workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
		TaskQueue:           queue,
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: ActivityRetryBackoffCoefficient,
			MaximumAttempts:    ActivityRetryMaximumAttempts,
			InitialInterval:    ActivityRetryInitialInterval,
			NonRetryableErrorTypes: []string{
				ErrNonRetryableUserExist,
			},
		},
	}
}

func validateRegisterUserInput(input *service.RegisterUserInput) error {
	if input == nil || input.User == nil {
		return entity.ErrEmptyUser()
	}
	return nil
}
