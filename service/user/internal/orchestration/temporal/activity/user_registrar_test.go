package activity_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	mock_activity "github.com/indrasaputra/arjuna/service/user/test/mock/orchestration/temporal/activity"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type RegisterUserActivitySuite struct {
	activity *activity.RegisterUserActivity

	conn *mock_activity.MockRegisterUserConnection
	db   *mock_activity.MockRegisterUserDatabase
}

func TestNewRegisterUserActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of RegisterUserActivity", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		assert.NotNil(t, st.activity)
	})
}

func TestRegisterUserActivity_CreateInAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user already exists", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.conn.EXPECT().CreateAccount(testCtx, user).Return(entity.ErrAlreadyExists())

		err := st.activity.CreateInAuth(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success create user", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.conn.EXPECT().CreateAccount(testCtx, user).Return(nil)

		err := st.activity.CreateInAuth(testCtx, user)

		assert.NoError(t, err)
	})
}

func TestRegisterUserActivity_HardDeleteInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("error when delete user from database", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(entity.ErrInternal(""))

		err := st.activity.HardDeleteInUser(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success delete user from conn", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := st.activity.HardDeleteInUser(testCtx, user.ID)

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
	co := mock_activity.NewMockRegisterUserConnection(ctrl)
	db := mock_activity.NewMockRegisterUserDatabase(ctrl)
	a := activity.NewRegisterUserActivity(co, db)
	return &RegisterUserActivitySuite{
		activity: a,
		conn:     co,
		db:       db,
	}
}
