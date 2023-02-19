package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type UserRegistrarSuite struct {
	registrar     *service.UserRegistrar
	orchestration *mock_service.MockRegisterUserOrchestration
}

func TestNewUserRegistrar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRegistrar", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		assert.NotNil(t, st.registrar)
	})
}

func TestUserRegistrar_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)

		id, err := st.registrar.Register(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Empty(t, id)
	})

	t.Run("name contains character other than alphabet", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		names := []string{
			"123",
			"First Us3r",
			"First User !!!",
			"F1rst User",
			"!@#$%^&*()",
		}

		for _, name := range names {
			user := &entity.User{Name: name}

			id, err := st.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidName(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("email is invalid", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		emails := []string{
			"@domain",
			"@domain.com",
			"domain.com",
		}

		for _, email := range emails {
			user := createTestUser()
			user.Email = email

			id, err := st.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("orchestration returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		st.orchestration.EXPECT().RegisterUser(testCtx, input).Return(nil, entity.ErrInternal("error"))

		id, err := st.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		input := &service.RegisterUserInput{User: user}
		st.orchestration.EXPECT().RegisterUser(testCtx, input).Return(&service.RegisterUserOutput{UserID: "user-id"}, nil)

		id, err := st.registrar.Register(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistrarSuite(ctrl *gomock.Controller) *UserRegistrarSuite {
	o := mock_service.NewMockRegisterUserOrchestration(ctrl)
	r := service.NewUserRegistrar(o)
	return &UserRegistrarSuite{
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
