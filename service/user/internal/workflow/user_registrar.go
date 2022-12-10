package workflow

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	tempoflow "go.temporal.io/sdk/workflow"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

const (
	// TaskQueueRegisterUser represents user registration.
	TaskQueueRegisterUser = "register-user"

	// ActivityTimeoutDefault sets to 2 seconds.
	ActivityTimeoutDefault = 2 * time.Second
	// ActivityKeycloakCreate is derived from struct name + method name. See activity registration in worker.
	ActivityKeycloakCreate = "KeycloakCreate"
	// ActivityKeycloakHardDelete is derived from struct name + method name. See activity registration in worker.
	ActivityKeycloakHardDelete = "KeycloakHardDelete"
	// ActivityPostgresInsert is derived from struct name + method name. See activity registration in worker.
	ActivityPostgresInsert = "PostgresInsert"
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
)

// RegisterUserExecutor is responsible to execute register user workflow.
type RegisterUserExecutor struct {
	client client.Client
}

// NewRegisterUserExecutor creates an instance of RegisterUserExecutor.
func NewRegisterUserExecutor(client client.Client) *RegisterUserExecutor {
	return &RegisterUserExecutor{client: client}
}

// RegisterUser runs the register users workflow.
func (r *RegisterUserExecutor) RegisterUser(ctx context.Context, input *service.RegisterUserInput) (*service.RegisterUserOutput, error) {
	opts := client.StartWorkflowOptions{
		ID:                 fmt.Sprintf("%s-%s", WorkflowNameRegisterUser, input.User.ID),
		TaskQueue:          TaskQueueRegisterUser,
		WorkflowRunTimeout: WorkflowTimeoutDefault,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: WorkflowRetryMaximumAttempts,
		},
	}
	wr, err := r.client.ExecuteWorkflow(ctx, opts, RegisterUser, input)
	if err != nil {
		return nil, err
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	var output *service.RegisterUserOutput
	err = wr.Get(ctx, &output)
	if err != nil {
		return nil, entity.ErrInternal("Something went wrong within our server. Please, try again")
	}
	return output, nil
}

// RegisterUser runs the user registration workflow.
func RegisterUser(ctx tempoflow.Context, input *service.RegisterUserInput) (*service.RegisterUserOutput, error) {
	if err := validateRegisterUserInput(input); err != nil {
		return nil, err
	}

	var id string
	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err := tempoflow.ExecuteActivity(ctx, ActivityKeycloakCreate, input.User).Get(ctx, &id)
	if err != nil {
		return nil, err
	}
	input.User.KeycloakID = id

	ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
	err = tempoflow.ExecuteActivity(ctx, ActivityPostgresInsert, input.User).Get(ctx, nil)
	if err != nil {
		ctx = createContextWithActivityOptions(ctx, ActivityTimeoutDefault, TaskQueueRegisterUser)
		_ = tempoflow.ExecuteActivity(ctx, ActivityKeycloakHardDelete, id).Get(ctx, nil)
		return nil, entity.ErrInternal("Something went wrong within our server. Please try again")
	}
	return &service.RegisterUserOutput{UserID: id}, nil
}

func createContextWithActivityOptions(tctx tempoflow.Context, timeout time.Duration, queue string) tempoflow.Context {
	opts := createActivityOptions(timeout, queue)
	return tempoflow.WithActivityOptions(tctx, opts)
}

func createActivityOptions(timeout time.Duration, queue string) tempoflow.ActivityOptions {
	return tempoflow.ActivityOptions{
		StartToCloseTimeout: timeout,
		TaskQueue:           queue,
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: ActivityRetryBackoffCoefficient,
			MaximumAttempts:    ActivityRetryMaximumAttempts,
			InitialInterval:    ActivityRetryInitialInterval,
		},
	}
}

func validateRegisterUserInput(input *service.RegisterUserInput) error {
	if input.User == nil {
		return entity.ErrEmptyUser()
	}
	return nil
}
