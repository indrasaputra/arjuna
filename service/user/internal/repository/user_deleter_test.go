package repository_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository"
	mock_repository "github.com/indrasaputra/arjuna/service/user/test/mock/repository"
)

var (
	testCtx = context.Background()
)

type UserDeleterExecutor struct {
	deleter  *repository.UserDeleter
	keycloak *mock_repository.MockDeleteUserRepository
	postgres *mock_repository.MockDeleteUserPostgres
}

func TestNewUserDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserDeleter", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		assert.NotNil(t, exec.deleter)
	})
}

func TestUserDeleter_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1"

	t.Run("get by id returns error", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.postgres.EXPECT().GetByID(testCtx, id).Return(nil, entity.ErrInternal("error"))

		user, err := exec.deleter.GetByID(testCtx, id)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("success get by id", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.postgres.EXPECT().GetByID(testCtx, id).Return(&entity.User{}, nil)

		user, err := exec.deleter.GetByID(testCtx, id)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestUserDeleter_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{}

	t.Run("keycloak returns error", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.keycloak.EXPECT().HardDelete(testCtx, user.KeycloakID).Return(entity.ErrInternal("error"))

		err := exec.deleter.HardDelete(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("postgres returns error", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.keycloak.EXPECT().HardDelete(testCtx, user.KeycloakID).Return(nil)
		exec.postgres.EXPECT().HardDelete(testCtx, user.ID).Return(entity.ErrInternal("error"))

		err := exec.deleter.HardDelete(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success delete user", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.keycloak.EXPECT().HardDelete(testCtx, user.KeycloakID).Return(nil)
		exec.postgres.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := exec.deleter.HardDelete(testCtx, user)

		assert.NoError(t, err)
	})
}

func createUserDeleterExecutor(ctrl *gomock.Controller) *UserDeleterExecutor {
	k := mock_repository.NewMockDeleteUserRepository(ctrl)
	p := mock_repository.NewMockDeleteUserPostgres(ctrl)
	d := repository.NewUserDeleter(k, p)
	return &UserDeleterExecutor{
		deleter:  d,
		keycloak: k,
		postgres: p,
	}
}
