package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type UserDeleterSuite struct {
	deleter *service.UserDeleter
	db      *mock_service.MockDeleteUserRepository
}

func TestNewUserDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserDeleter", func(t *testing.T) {
		st := createUserDeleterSuite(ctrl)
		assert.NotNil(t, st.deleter)
	})
}

func TestUserDeleter_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("delete from db returns error", func(t *testing.T) {
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st := createUserDeleterSuite(ctrl)
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(errReturn)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("success hard delete user", func(t *testing.T) {
		user := createTestUser()

		st := createUserDeleterSuite(ctrl)
		st.db.EXPECT().HardDelete(testCtx, user.ID).Return(nil)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func createUserDeleterSuite(ctrl *gomock.Controller) *UserDeleterSuite {
	db := mock_service.NewMockDeleteUserRepository(ctrl)
	d := service.NewUserDeleter(db)
	return &UserDeleterSuite{
		deleter: d,
		db:      db,
	}
}
