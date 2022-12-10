package workflow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/activity"
	tempomock "go.temporal.io/sdk/mocks"
	"go.temporal.io/sdk/testsuite"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	"github.com/indrasaputra/arjuna/service/user/internal/workflow"
)

var (
	testCtx = context.Background()
)

type RegisterUserExecutorSuite struct {
	exec   *workflow.RegisterUserExecutor
	client *tempomock.Client
}

func TestNewRegisterUserExecutor(t *testing.T) {
	t.Run("successfully create an instance of RegisterUserExecutor", func(t *testing.T) {
		st := createRegisterUserExecutorSuite()
		assert.NotNil(t, st.exec)
	})
}

func TestRegisterUserExecutor_RegisterUser(t *testing.T) {
	t.Run("execute workflow returns error", func(t *testing.T) {
		st := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(nil, errors.New("temporal is down"))

		res, err := st.exec.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow run returns error", func(t *testing.T) {
		st := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(errors.New("workflow run error"))

		res, err := st.exec.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow is executed successfully", func(t *testing.T) {
		st := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		st.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(nil)

		res, err := st.exec.RegisterUser(testCtx, input)

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

type RegisterUserSuite struct {
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
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

	t.Run("KeycloakCreate activity returns error", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()

		st.env.OnActivity(workflow.ActivityKeycloakCreate, mock.Anything, input.User).Return("", errors.New("keycloak error"))

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("KeycloakCreate activity returns error", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()

		st.env.OnActivity(workflow.ActivityKeycloakCreate, mock.Anything, input.User).Return("", errors.New("keycloak error"))

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("PostgresInsert activity returns error", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()
		id := "1"

		st.env.OnActivity(workflow.ActivityKeycloakCreate, mock.Anything, mock.Anything).Return(id, nil)
		st.env.OnActivity(workflow.ActivityPostgresInsert, mock.Anything, mock.Anything).Return(errors.New("keycloak error"))
		st.env.OnActivity(workflow.ActivityKeycloakHardDelete, mock.Anything, id).Return(nil)

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.Error(t, st.env.GetWorkflowError())
	})

	t.Run("workflow is executed successfully", func(t *testing.T) {
		st := createRegisterUserSuite()
		input := createRegisterUserInput()
		id := "1"

		st.env.OnActivity(workflow.ActivityKeycloakCreate, mock.Anything, mock.Anything).Return(id, nil)
		st.env.OnActivity(workflow.ActivityPostgresInsert, mock.Anything, mock.Anything).Return(nil)

		st.env.ExecuteWorkflow(workflow.RegisterUser, input)

		assert.True(t, st.env.IsWorkflowCompleted())
		assert.NoError(t, st.env.GetWorkflowError())

		var res *service.RegisterUserOutput
		_ = st.env.GetWorkflowResult(&res)
		assert.Equal(t, id, res.UserID)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createRegisterUserInput() *service.RegisterUserInput {
	return &service.RegisterUserInput{
		User: createTestUser(),
	}
}

func createRegisterUserExecutorSuite() *RegisterUserExecutorSuite {
	c := &tempomock.Client{}
	e := workflow.NewRegisterUserExecutor(c)
	return &RegisterUserExecutorSuite{
		exec:   e,
		client: c,
	}
}

func createRegisterUserSuite() *RegisterUserSuite {
	s := &RegisterUserSuite{}
	s.env = s.NewTestWorkflowEnvironment()

	kc := &keycloak.User{}
	pg := &postgres.User{}
	s.env.RegisterActivityWithOptions(kc, activity.RegisterOptions{Name: "Keycloak", SkipInvalidStructFunctions: true})
	s.env.RegisterActivityWithOptions(pg, activity.RegisterOptions{Name: "Postgres", SkipInvalidStructFunctions: true})

	return s
}
