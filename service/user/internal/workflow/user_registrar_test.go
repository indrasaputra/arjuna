package workflow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tempomock "go.temporal.io/sdk/mocks"

	"github.com/indrasaputra/arjuna/service/user/entity"
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
		suite := createRegisterUserExecutorSuite()
		assert.NotNil(t, suite.exec)
	})
}

func TestRegisterUserExecutor_RegisterUser(t *testing.T) {
	t.Run("execute workflow returns error", func(t *testing.T) {
		suite := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}

		suite.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(nil, errors.New("temporal is down"))

		res, err := suite.exec.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow run returns error", func(t *testing.T) {
		suite := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		suite.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(errors.New("workflow run error"))

		res, err := suite.exec.RegisterUser(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("workflow is executed successfully", func(t *testing.T) {
		suite := createRegisterUserExecutorSuite()
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		wr := &tempomock.WorkflowRun{}

		suite.client.
			On("ExecuteWorkflow", testCtx, mock.Anything, mock.AnythingOfType("func(internal.Context, *service.RegisterUserInput) (*service.RegisterUserOutput, error)"), input).
			Return(wr, nil)
		wr.On("GetID").Return("")
		wr.On("GetRunID").Return("")
		wr.On("Get", testCtx, mock.Anything).Return(nil)

		res, err := suite.exec.RegisterUser(testCtx, input)

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		Name:  "First User",
		Email: "first@user.com",
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
