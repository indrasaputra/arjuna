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

type RegisterUserActivitySuite struct {
	activity *activity.RegisterUserActivity

	vendor *mock_activity.MockRegisterUserVendor
	db     *mock_activity.MockRegisterUserDatabase
}

func TestNewRegisterUserActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of RegisterUserActivity", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		assert.NotNil(t, st.activity)
	})
}

func TestRegisterUserActivity_CreateInKeycloak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user already exists", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.vendor.EXPECT().Create(testCtx, user).Return("", entity.ErrAlreadyExists())

		id, err := st.activity.CreateInKeycloak(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create user", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.vendor.EXPECT().Create(testCtx, user).Return("1", nil)

		id, err := st.activity.CreateInKeycloak(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestRegisterUserActivity_HardDeleteFromKeycloak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error when delete user from vendor", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.vendor.EXPECT().HardDelete(testCtx, user.ID).Return(entity.ErrInternal(""))

		err := st.activity.HardDeleteFromKeycloak(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success delete user from vendor", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.vendor.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := st.activity.HardDeleteFromKeycloak(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func TestRegisterUserActivity_UpdateKeycloakID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("keycloak id already exists in database", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().UpdateKeycloakID(testCtx, user.ID, user.KeycloakID).Return(entity.ErrAlreadyExists())

		err := st.activity.UpdateKeycloakID(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success update keycloak id to database", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().UpdateKeycloakID(testCtx, user.ID, user.KeycloakID).Return(nil)

		err := st.activity.UpdateKeycloakID(testCtx, user)

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

func createRegisterUserActivitySuite(ctrl *gomock.Controller) *RegisterUserActivitySuite {
	kc := mock_activity.NewMockRegisterUserVendor(ctrl)
	db := mock_activity.NewMockRegisterUserDatabase(ctrl)
	a := activity.NewRegisterUserActivity(kc, db)
	return &RegisterUserActivitySuite{
		activity: a,
		vendor:   kc,
		db:       db,
	}
}
