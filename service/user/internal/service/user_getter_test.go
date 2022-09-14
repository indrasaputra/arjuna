package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type UserGetterExecutor struct {
	getter *service.UserGetter
	repo   *mock_service.MockGetUserRepository
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

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx, service.DefaultGetAllUsersLimit).Return([]*entity.User{}, entity.ErrInternal(""))

		res, err := exec.getter.GetAll(testCtx, service.DefaultGetAllUsersLimit)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Empty(t, res)
	})

	t.Run("repository returns empty list", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx, service.DefaultGetAllUsersLimit).Return([]*entity.User{}, nil)

		res, err := exec.getter.GetAll(testCtx, service.DefaultGetAllUsersLimit)

		assert.Nil(t, err)
		assert.Empty(t, res)
	})

	t.Run("successfully all users", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx, service.DefaultGetAllUsersLimit).Return([]*entity.User{{}}, nil)

		res, err := exec.getter.GetAll(testCtx, service.DefaultGetAllUsersLimit)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("successfully all users with limit = 0", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx, service.DefaultGetAllUsersLimit).Return([]*entity.User{{}}, nil)

		res, err := exec.getter.GetAll(testCtx, 0)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("successfully all users with limit > 10", func(t *testing.T) {
		exec := createUserGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx, service.DefaultGetAllUsersLimit).Return([]*entity.User{{}}, nil)

		res, err := exec.getter.GetAll(testCtx, 100)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})
}

func createUserGetterExecutor(ctrl *gomock.Controller) *UserGetterExecutor {
	repo := mock_service.NewMockGetUserRepository(ctrl)
	reg := service.NewUserGetter(repo)
	return &UserGetterExecutor{
		getter: reg,
		repo:   repo,
	}
}
