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

type UserRegistrarExecutor struct {
	registrar     *service.UserRegistrar
	orchestration *mock_service.MockRegisterUserOrchestration
}

func TestNewUserRegistrar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRegistrar", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)
		assert.NotNil(t, exec.registrar)
	})
}

func TestUserRegistrar_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty user is prohibited", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)

		id, err := exec.registrar.Register(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Empty(t, id)
	})

	t.Run("name contains character other than alphabet", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)
		names := []string{
			"123",
			"First Us3r",
			"First User !!!",
			"F1rst User",
			"!@#$%^&*()",
		}

		for _, name := range names {
			user := &entity.User{Name: name}

			id, err := exec.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidName(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("email is invalid", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)
		emails := []string{
			"@domain",
			"@domain.com",
			"domain.com",
		}

		for _, email := range emails {
			user := createTestUser()
			user.Email = email

			id, err := exec.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("orchestration returns error", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		exec.orchestration.EXPECT().RegisterUser(testCtx, input).Return(nil, entity.ErrInternal("error"))

		id, err := exec.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		exec := createUserRegistrarExecutor(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		exec.orchestration.EXPECT().RegisterUser(testCtx, input).Return(&service.RegisterUserOutput{UserID: "user-id"}, nil)

		id, err := exec.registrar.Register(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistrarExecutor(ctrl *gomock.Controller) *UserRegistrarExecutor {
	o := mock_service.NewMockRegisterUserOrchestration(ctrl)
	r := service.NewUserRegistrar(o)
	return &UserRegistrarExecutor{
		registrar:     r,
		orchestration: o,
	}
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:         "1",
		KeycloakID: "1",
		Name:       "First User",
		Email:      "first@user.com",
	}
}
