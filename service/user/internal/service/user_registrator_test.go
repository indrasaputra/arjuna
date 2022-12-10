package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

var (
	testCtx = context.Background()
)

type UserRegistratorExecutor struct {
	registrator *service.UserRegistrator
	workflow    *mock_service.MockRegisterUserWorkflow
}

func TestNewUserRegistrator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRegistrator", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		assert.NotNil(t, exec.registrator)
	})
}

func TestUserRegistrator_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty user is prohibited", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)

		id, err := exec.registrator.Register(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Empty(t, id)
	})

	t.Run("name contains character other than alphabet", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		names := []string{
			"123",
			"First Us3r",
			"First User !!!",
			"F1rst User",
			"!@#$%^&*()",
		}

		for _, name := range names {
			user := &entity.User{Name: name}

			id, err := exec.registrator.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidName(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("email is invalid", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		emails := []string{
			"@domain",
			"@domain.com",
			"domain.com",
		}

		for _, email := range emails {
			user := createTestUser()
			user.Email = email

			id, err := exec.registrator.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("workflow returns error", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		exec.workflow.EXPECT().RegisterUser(testCtx, input).Return(nil, entity.ErrInternal("error"))

		id, err := exec.registrator.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		exec.workflow.EXPECT().RegisterUser(testCtx, input).Return(&service.RegisterUserOutput{UserID: "user-id"}, nil)

		id, err := exec.registrator.Register(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistratorExecutor(ctrl *gomock.Controller) *UserRegistratorExecutor {
	r := mock_service.NewMockRegisterUserWorkflow(ctrl)
	rg := service.NewUserRegistrator(r)
	return &UserRegistratorExecutor{
		registrator: rg,
		workflow:    r,
	}
}

func createTestUser() *entity.User {
	return &entity.User{
		Name:  "First User",
		Email: "first@user.com",
	}
}
