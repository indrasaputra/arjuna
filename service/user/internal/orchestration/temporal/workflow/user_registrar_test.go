package workflow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/activity"
	tempomock "go.temporal.io/sdk/mocks"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/connection/auth"
	"github.com/indrasaputra/arjuna/service/user/internal/connection/wallet"
	orcact "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type RegisterUserWorkflowSuite struct {
	workflow *workflow.RegisterUserWorkflow
	client   *tempomock.Client
}

func TestNewRegisterUserWorkflow(t *testing.T) {
	t.Run("successfully create an instance of RegisterUserWorkflow", func(t *testing.T) {
		st := createRegisterUserWorkflowSuite()
		assert.NotNil(t, st.workflow)
	})
}

func TestRegisterUserWorkflow_RegisterUser(t *testing.T) {
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("execute workflow returns error", func(t *testing.T) {
		st := createRegisterUserWorkflowSuite()
		user := createTestUser()
		input := &entity.RegisterUserInput{User: user}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *entity.RegisterUserInput) (*entity.RegisterUserOutput, error)"), input).
			Return(nil, errors.New("temporal is down"))

		res, err := st.workflow.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow run returns user already exists error", func(t *testing.T) {
		st := createRegisterUserWorkflowSuite()
		user := createTestUser()
		input := &entity.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *entity.RegisterUserInput) (*entity.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(temporal.NewNonRetryableApplicationError("", workflow.ErrNonRetryableUserExist, errors.New("")))

		res, err := st.workflow.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow run returns internal error", func(t *testing.T) {
		st := createRegisterUserWorkflowSuite()
		user := createTestUser()
		input := &entity.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *entity.RegisterUserInput) (*entity.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(errors.New("workflow run error"))

		res, err := st.workflow.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow is executed successfully", func(t *testing.T) {
		st := createRegisterUserWorkflowSuite()
		user := createTestUser()
		input := &entity.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *entity.RegisterUserInput) (*entity.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(nil)

		res, err := st.workflow.RegisterUser(testCtx, input)

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

type RegisterUserSuite struct {
	env *testsuite.TestWorkflowEnvironment
	testsuite.WorkflowTestSuite
}

func TestRegisterUser(t *testing.T) {
	t.Run("input is invalid", func(t *testing.T) {
		st := createRegisterUserSuite()

		st.env.ExecuteWorkflow(workflow.RegisterUser, nil)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("input doesn't have user struct", func(t *testing.T) {
		st := createRegisterUserSuite()

		input := createRegisterUserInput()
		input.User = nil
		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("AuthCreate activity returns error", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()

		st.env.OnActivity(workflow.ActivityAuthCreate, mock.Anything, input.User).Return(errors.New("auth error"))

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("WalletCreate activity returns error", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()

		st.env.OnActivity(workflow.ActivityAuthCreate, mock.Anything, input.User).Return(nil)
		st.env.OnActivity(workflow.ActivityWalletCreate, mock.Anything, input.User).Return(errors.New("wallet error"))
		st.env.OnActivity(workflow.ActivityUserHardDelete, mock.Anything, input.User.ID).Return(nil)

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("workflow is executed successfully", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()

		st.env.OnActivity(workflow.ActivityAuthCreate, mock.Anything, input.User).Return(nil)
		st.env.OnActivity(workflow.ActivityWalletCreate, mock.Anything, input.User).Return(nil)

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.NoError(t, st.env.GetWorkflowError())

		var res *entity.RegisterUserOutput
		_ = st.env.GetWorkflowResult(&res)
		assert.NotNil(t, res)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    uuid.Must(uuid.NewV7()),
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createRegisterUserInput() *entity.RegisterUserInput {
	return &entity.RegisterUserInput{
		User: createTestUser(),
	}
}

func createRegisterUserWorkflowSuite() *RegisterUserWorkflowSuite {
	c := &tempomock.Client{}
	w := workflow.NewRegisterUserWorkflow(c)
	return &RegisterUserWorkflowSuite{
		workflow: w,
		client:   c,
	}
}

func createRegisterUserSuite() *RegisterUserSuite {
	s := &RegisterUserSuite{}
	s.env = s.NewTestWorkflowEnvironment()

	at := &auth.Auth{}
	wt := &wallet.Wallet{}
	pg := &postgres.User{}
	uc := orcact.NewRegisterUserActivity(at, wt, pg)

	s.env.RegisterActivityWithOptions(uc, activity.RegisterOptions{Name: "RegisterUserActivity", SkipInvalidStructFunctions: true})

	return s
}
