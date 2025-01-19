package workflow

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	tempflow "go.temporal.io/sdk/workflow"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	// TaskQueueRegisterUser represents user registration.
	TaskQueueRegisterUser = "register-user"

	// ActivityTimeoutDefault sets to 2 seconds.
	ActivityTimeoutDefault = 2 * time.Second
	// ActivityAuthCreate is derived from struct name + method name. See activity registration in worker.
	ActivityAuthCreate = "RegisterUserActivityCreateAccount"
	// ActivityUserHardDelete is derived from struct name + method name. See activity registration in worker.
	ActivityUserHardDelete = "RegisterUserActivityHardDeleteInUser"
	// ActivityWalletCreate is derived from struct name + method name. See activity registration in worker.
	ActivityWalletCreate = "RegisterUserActivityCreateWallet"
	// ActivityRetryBackoffCoefficient sets to 2.
	ActivityRetryBackoffCoefficient = 2
	// ActivityRetryMaximumAttempts sets to 3.
	ActivityRetryMaximumAttempts = 1
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
func (r *RegisterUserWorkflow) RegisterUser(ctx context.Context, input *entity.RegisterUserInput) (*entity.RegisterUserOutput, error) {
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
		slog.ErrorContext(ctx, "[RegisterUserWorkflow-RegisterUser] fail to start workflow", "error", err)
		return nil, entity.ErrInternal("Something went wrong within our server. Please, try again")
	}
	slog.InfoContext(ctx, "[RegisterUserWorkflow-RegisterUser] started workflow", "workflow-id", wr.GetID(), "run-id", wr.GetRunID())

	var output *entity.RegisterUserOutput
	err = wr.Get(ctx, &output)
	if err != nil {
		var appErr *temporal.ApplicationError
		if errors.As(err, &appErr) && appErr.Type() == ErrNonRetryableUserExist {
			return nil, entity.ErrAlreadyExists()
		}
		slog.ErrorContext(ctx, "[RegisterUserWorkflow-RegisterUser] error get workflow result", "error", err)
		return nil, entity.ErrInternal("Something went wrong within our server. Please, try again")
	}
	return output, nil
}

// RegisterUser runs the user registration workflow.
func RegisterUser(ctx tempflow.Context, input *entity.RegisterUserInput) (*entity.RegisterUserOutput, error) {
	if err := validateRegisterUserInput(input); err != nil {
		return nil, err
	}

	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err := tempflow.ExecuteActivity(ctx, ActivityAuthCreate, input.User).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err = tempflow.ExecuteActivity(ctx, ActivityWalletCreate, input.User).Get(ctx, nil)
	if err != nil {
		ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
		_ = tempflow.ExecuteActivity(ctx, ActivityUserHardDelete, input.User.ID).Get(ctx, nil)
		return nil, entity.ErrInternal("Something went wrong within our server. Please try again")
	}
	return &entity.RegisterUserOutput{}, nil
}

func createContextWithActivityOptions(tempoCtx tempflow.Context, timeout time.Duration, queue string) tempflow.Context {
	opts := createActivityOptions(timeout, queue)
	return tempflow.WithActivityOptions(tempoCtx, opts)
}

func createActivityOptions(timeout time.Duration, queue string) tempflow.ActivityOptions {
	return tempflow.ActivityOptions{
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

func validateRegisterUserInput(input *entity.RegisterUserInput) error {
	if input == nil || input.User == nil {
		return entity.ErrEmptyUser()
	}
	return nil
}
