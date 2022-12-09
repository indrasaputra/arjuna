package workflow

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	tempoflow "go.temporal.io/sdk/workflow"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

const (
	// TaskQueueRegisterUser represents user registration.
	TaskQueueRegisterUser = "register-user"
	// TimeoutDefault sets to 5 seconds.
	TimeoutDefault = 5 * time.Second
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
		ID:        fmt.Sprintf("register-user-%s", input.User.ID),
		TaskQueue: TaskQueueRegisterUser,
	}
	wr, err := r.client.ExecuteWorkflow(ctx, opts, RegisterUser, input)
	if err != nil {
		return nil, err
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	var output *service.RegisterUserOutput
	err = wr.Get(ctx, &output)
	if err != nil {
		return nil, err
	}
	return output, err
}

// RegisterUser runs the user registration workflow.
func RegisterUser(ctx tempoflow.Context, input *service.RegisterUserInput) (*service.RegisterUserOutput, error) {
	if err := validateRegisterUserInput(input); err != nil {
		return nil, err
	}

	var id string
	options := createActivityOptions(TimeoutDefault, TaskQueueRegisterUser)
	ctx = tempoflow.WithActivityOptions(ctx, options)
	err := tempoflow.ExecuteActivity(ctx, "KeycloakCreate", input.User).Get(ctx, &id)
	if err != nil {
		return nil, err
	}
	input.User.KeycloakID = id

	options = createActivityOptions(TimeoutDefault, TaskQueueRegisterUser)
	ctx = tempoflow.WithActivityOptions(ctx, options)
	err = tempoflow.ExecuteActivity(ctx, "PostgresInsert", input.User).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &service.RegisterUserOutput{UserID: id}, nil
}

func createActivityOptions(timeout time.Duration, queue string) tempoflow.ActivityOptions {
	return tempoflow.ActivityOptions{
		StartToCloseTimeout: timeout,
		TaskQueue:           queue,
	}
}

func validateRegisterUserInput(input *service.RegisterUserInput) error {
	if input.User == nil {
		return entity.ErrEmptyUser()
	}
	return nil
}
