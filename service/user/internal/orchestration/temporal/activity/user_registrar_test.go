package activity_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	mock_activity "github.com/indrasaputra/arjuna/service/user/test/mock/orchestration/temporal/activity"
)

var (
	testCtx = context.Background()
)

type RegisterUserActivitySuite struct {
	activity *activity.RegisterUserActivity

	auth   *mock_activity.MockRegisterUserAuthConnection
	wallet *mock_activity.MockRegisterUserWalletConnection
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

func TestRegisterUserActivity_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("auth already exists", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.auth.EXPECT().CreateAccount(testCtx, user).Return(entity.ErrAlreadyExists())

		err := st.activity.CreateAccount(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success create account", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.auth.EXPECT().CreateAccount(testCtx, user).Return(nil)

		err := st.activity.CreateAccount(testCtx, user)

		assert.NoError(t, err)
	})
}

func TestRegisterUserActivity_CreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("wallet already exists", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.wallet.EXPECT().CreateWallet(testCtx, user).Return(entity.ErrAlreadyExists())

		err := st.activity.CreateWallet(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success create wallet", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.wallet.EXPECT().CreateWallet(testCtx, user).Return(nil)

		err := st.activity.CreateWallet(testCtx, user)

		assert.NoError(t, err)
	})
}

func TestRegisterUserActivity_HardDeleteInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error when delete user from database", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(assert.AnError)

		err := st.activity.HardDeleteInUser(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success delete user from auth", func(t *testing.T) {
		st := createRegisterUserActivitySuite(ctrl)
		user := createTestUser()
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := st.activity.HardDeleteInUser(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    uuid.Must(uuid.NewV7()),
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createRegisterUserActivitySuite(ctrl *gomock.Controller) *RegisterUserActivitySuite {
	ac := mock_activity.NewMockRegisterUserAuthConnection(ctrl)
	wc := mock_activity.NewMockRegisterUserWalletConnection(ctrl)
	db := mock_activity.NewMockRegisterUserDatabase(ctrl)
	a := activity.NewRegisterUserActivity(ac, wc, db)
	return &RegisterUserActivitySuite{
		activity: a,
		auth:     ac,
		wallet:   wc,
		db:       db,
	}
}
