package repository_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

var (
	testCtx = context.Background()
)

type UserRegistratorExecutor struct {
	registrator *repository.UserRegistrator
	keycloak    *mock_service.MockRegisterUserRepository
	postgres    *mock_service.MockRegisterUserRepository
}

func TestNewUserRegistrator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRegistrator", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		assert.NotNil(t, exec.registrator)
	})
}

func TestUserRegistrator_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user is empty", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)

		err := exec.registrator.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	user := &entity.User{}

	t.Run("keycloak returns error", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		exec.keycloak.EXPECT().Insert(testCtx, user).Return(entity.ErrInternal("error"))

		err := exec.registrator.Insert(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("postgres returns error", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		exec.keycloak.EXPECT().Insert(testCtx, user).Return(nil)
		exec.postgres.EXPECT().Insert(testCtx, user).Return(entity.ErrInternal("error"))

		err := exec.registrator.Insert(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success insert user", func(t *testing.T) {
		exec := createUserRegistratorExecutor(ctrl)
		exec.keycloak.EXPECT().Insert(testCtx, user).Return(nil)
		exec.postgres.EXPECT().Insert(testCtx, user).Return(nil)

		err := exec.registrator.Insert(testCtx, user)

		assert.NoError(t, err)
	})
}

func createUserRegistratorExecutor(ctrl *gomock.Controller) *UserRegistratorExecutor {
	k := mock_service.NewMockRegisterUserRepository(ctrl)
	p := mock_service.NewMockRegisterUserRepository(ctrl)
	reg := repository.NewUserRegistrator(k, p)
	return &UserRegistratorExecutor{
		registrator: reg,
		keycloak:    k,
		postgres:    p,
	}
}
