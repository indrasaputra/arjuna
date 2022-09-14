package repository_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type UserGetterExecutor struct {
	getter   *repository.UserGetter
	keycloak *mock_service.MockGetUserRepository
	postgres *mock_service.MockGetUserRepository
}

func TestNewUserGetter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserGetter", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		assert.NotNil(t, exec.getter)
	})
}

func TestUserGetter_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("keycloak returns error", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.keycloak.EXPECT().GetAll(testCtx).Return(nil, entity.ErrInternal("error"))

		users, err := exec.getter.GetAll(testCtx)

		assert.Error(t, err)
		assert.Empty(t, users)
	})

	t.Run("postgres returns error", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.keycloak.EXPECT().GetAll(testCtx).Return([]*entity.User{{}}, nil)
		exec.postgres.EXPECT().GetAll(testCtx).Return(nil, entity.ErrInternal("error"))

		users, err := exec.getter.GetAll(testCtx)

		assert.Error(t, err)
		assert.Empty(t, users)
	})

	t.Run("success get all users", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.keycloak.EXPECT().GetAll(testCtx).Return([]*entity.User{{}}, nil)
		exec.postgres.EXPECT().GetAll(testCtx).Return([]*entity.User{{}}, nil)

		users, err := exec.getter.GetAll(testCtx)

		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})
}

func createUserGetterExecutor(ctrl *gomock.Controller) *UserGetterExecutor {
	k := mock_service.NewMockGetUserRepository(ctrl)
	p := mock_service.NewMockGetUserRepository(ctrl)
	g := repository.NewUserGetter(k, p)
	return &UserGetterExecutor{
		getter:   g,
		keycloak: k,
		postgres: p,
	}
}
