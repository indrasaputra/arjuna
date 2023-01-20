package activity_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	mock_activity "github.com/indrasaputra/arjuna/service/user/test/mock/orchestration/temporal/activity"
)

var (
	testCtx = context.Background()
)

type RegisterUserActivityExecutor struct {
	activity *activity.RegisterUserActivity

	vendor *mock_activity.MockRegisterUserVendor
	db     *mock_activity.MockRegisterUserDatabase
}

func TestNewRegisterUserActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of RegisterUserActivity", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		assert.NotNil(t, exec.activity)
	})
}

func TestRegisterUserActivity_CreateInKeycloak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user already exists", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.vendor.EXPECT().Create(testCtx, user).Return("", entity.ErrAlreadyExists())

		id, err := exec.activity.CreateInKeycloak(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create user", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.vendor.EXPECT().Create(testCtx, user).Return("1", nil)

		id, err := exec.activity.CreateInKeycloak(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestRegisterUserActivity_HardDeleteFromKeycloak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error when delete user from vendor", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.vendor.EXPECT().HardDelete(testCtx, user.ID).Return(entity.ErrInternal(""))

		err := exec.activity.HardDeleteFromKeycloak(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success delete user from vendor", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.vendor.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := exec.activity.HardDeleteFromKeycloak(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func TestRegisterUserActivity_InsertToDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user already exists in database", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.db.EXPECT().Insert(testCtx, user).Return(entity.ErrAlreadyExists())

		err := exec.activity.InsertToDatabase(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success insert user to database", func(t *testing.T) {
		exec := createRegisterUserActivityExecutor(ctrl)
		user := createTestUser()
		exec.db.EXPECT().Insert(testCtx, user).Return(nil)

		err := exec.activity.InsertToDatabase(testCtx, user)

		assert.NoError(t, err)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    "first-id",
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createRegisterUserActivityExecutor(ctrl *gomock.Controller) *RegisterUserActivityExecutor {
	kc := mock_activity.NewMockRegisterUserVendor(ctrl)
	db := mock_activity.NewMockRegisterUserDatabase(ctrl)
	a := activity.NewRegisterUserActivity(kc, db)
	return &RegisterUserActivityExecutor{
		activity: a,
		vendor:   kc,
		db:       db,
	}
}
