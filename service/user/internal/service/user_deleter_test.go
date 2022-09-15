package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

var (
	testUserID = "1"
)

type UserDeleterExecutor struct {
	deleter *service.UserDeleter
	repo    *mock_service.MockDeleteUserRepository
}

func TestNewUserDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserDeleter", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		assert.NotNil(t, exec.deleter)
	})
}

func TestUserDeleter_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("get user by email returns error", func(t *testing.T) {
		exec := createUserDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByID(testCtx, testUserID).Return(nil, entity.ErrInternal(""))

		err := exec.deleter.HardDelete(testCtx, testUserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("repo returns error when delete", func(t *testing.T) {
		user := &entity.User{}
		exec := createUserDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByID(testCtx, testUserID).Return(user, nil)
		exec.repo.EXPECT().HardDelete(testCtx, user).Return(entity.ErrInternal(""))

		err := exec.deleter.HardDelete(testCtx, testUserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("success delete user", func(t *testing.T) {
		user := &entity.User{}
		exec := createUserDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByID(testCtx, testUserID).Return(user, nil)
		exec.repo.EXPECT().HardDelete(testCtx, user).Return(nil)

		err := exec.deleter.HardDelete(testCtx, testUserID)

		assert.NoError(t, err)
	})
}

func createUserDeleterExecutor(ctrl *gomock.Controller) *UserDeleterExecutor {
	r := mock_service.NewMockDeleteUserRepository(ctrl)
	d := service.NewUserDeleter(r)
	return &UserDeleterExecutor{
		deleter: d,
		repo:    r,
	}
}
