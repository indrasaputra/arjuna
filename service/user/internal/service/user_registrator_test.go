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
	repo        *mock_service.MockRegisterUserRepository
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
			"Zlatan 1brahimovic",
			"Zlatan Ibrahimovic !!!",
			"5latan Ibrahimovic",
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
			user := &entity.User{
				Name:  "Zlatan Ibrahimovic",
				Email: email,
			}

			id, err := exec.registrator.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("repo returns error", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		user := &entity.User{
			Name:  "Zlatan Ibrahimovic",
			Email: "zlatan@ibrahimovic.com",
		}
		exec.repo.EXPECT().Insert(testCtx, user).Return(entity.ErrInternal("error"))

		id, err := exec.registrator.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		user := &entity.User{
			Name:  "Zlatan Ibrahimovic",
			Email: "zlatan@ibrahimovic.com",
		}
		exec.repo.EXPECT().Insert(testCtx, user).Return(nil)

		id, err := exec.registrator.Register(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistratorExecutor(ctrl *gomock.Controller) *UserRegistratorExecutor {
	r := mock_service.NewMockRegisterUserRepository(ctrl)
	rg := service.NewUserRegistrator(r)
	return &UserRegistratorExecutor{
		registrator: rg,
		repo:        r,
	}
}
