package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type UserDeleterSuite struct {
	deleter  *service.UserDeleter
	database *mock_service.MockDeleteUserRepository
	unit     *mock_uow.MockUnitOfWork
	tx       *mock_uow.MockTx
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

	t.Run("get user by id returns error", func(t *testing.T) {
		user := createTestUser()
		errReturn := entity.ErrInternal("")
		st := createUserDeleterSuite(ctrl)
		st.database.EXPECT().GetByID(testCtx, user.ID).Return(nil, errReturn)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("unit of work begin returns error", func(t *testing.T) {
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st := createUserDeleterSuite(ctrl)
		st.database.EXPECT().GetByID(testCtx, user.ID).Return(user, nil)
		st.unit.EXPECT().Begin(testCtx).Return(nil, errReturn)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("delete from database returns error and rollback", func(t *testing.T) {
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st := createUserDeleterSuite(ctrl)
		st.database.EXPECT().GetByID(testCtx, user.ID).Return(user, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.database.EXPECT().HardDeleteWithTx(testCtx, st.tx, user.ID).Return(errReturn)
		st.unit.EXPECT().Finish(testCtx, st.tx, errReturn).Return(errReturn)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("finish returns error", func(t *testing.T) {
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st := createUserDeleterSuite(ctrl)
		st.database.EXPECT().GetByID(testCtx, user.ID).Return(user, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.database.EXPECT().HardDeleteWithTx(testCtx, st.tx, user.ID).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(errReturn)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("success hard delete user", func(t *testing.T) {
		user := createTestUser()

		st := createUserDeleterSuite(ctrl)
		st.database.EXPECT().GetByID(testCtx, user.ID).Return(user, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.database.EXPECT().HardDeleteWithTx(testCtx, st.tx, user.ID).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(nil)

		err := st.deleter.HardDelete(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func createUserDeleterSuite(ctrl *gomock.Controller) *UserDeleterSuite {
	db := mock_service.NewMockDeleteUserRepository(ctrl)
	u := mock_uow.NewMockUnitOfWork(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	d := service.NewUserDeleter(u, db)
	return &UserDeleterSuite{
		deleter:  d,
		database: db,
		unit:     u,
		tx:       tx,
	}
}
