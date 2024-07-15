package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

const (
	queryUpdateRecordStatus = "UPDATE users_outbox SET status = ? WHERE id = ?"
)

type UserOutboxSuite struct {
	outbox *postgres.UserOutbox
	db     *mock_uow.MockDB
	tx     *mock_uow.MockTx
}

func TestNewUserOutbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of User", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		assert.NotNil(t, st.outbox)
	})
}

func TestUserOutbox_InsertWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := "INSERT INTO " +
		"users_outbox (id, status, payload, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil tx is prohibited", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)

		err := st.outbox.InsertWithTx(testCtx, nil, nil)

		assert.Error(t, err)
	})

	t.Run("nil outbox is prohibited", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)

		err := st.outbox.InsertWithTx(testCtx, st.tx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("nil payload is prohibited", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		out := createTestUserOutbox()
		out.Payload = nil

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("insert duplicate outbox", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		out := createTestUserOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		out := createTestUserOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		out := createTestUserOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), nil)

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_GetAllReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "SELECT id, status, payload FROM users_outbox WHERE status = ? ORDER BY created_at ASC LIMIT ? FOR UPDATE"
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, entity.UserOutboxStatusReady, limit).
			Return(errPostgresInternal)

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, entity.UserOutboxStatusReady, limit).
			Return(nil)

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestUserOutbox_SetProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set processed returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusProcessed, out.ID).
			Return(int64(0), errPostgresInternal)

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set processed success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusProcessed, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetDelivered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set delivered returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusDelivered, out.ID).
			Return(int64(0), errPostgresInternal)

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set delivered success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusDelivered, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set failed returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusFailed, out.ID).
			Return(int64(0), errPostgresInternal)

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set failed success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.UserOutboxStatusFailed, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetRecordStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set status returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		status := entity.UserOutboxStatusProcessed
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, status, out.ID).
			Return(int64(0), errPostgresInternal)

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.Error(t, err)
	})

	t.Run("set status success", func(t *testing.T) {
		out := createTestUser()
		status := entity.UserOutboxStatusProcessed
		st := createUserOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, status, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.NoError(t, err)
	})
}

func createTestUserOutbox() *entity.UserOutbox {
	user := createTestUser()
	return &entity.UserOutbox{
		ID:      "1",
		Status:  entity.UserOutboxStatusReady,
		Payload: user,
	}
}

func createUserOutboxSuite(ctrl *gomock.Controller) *UserOutboxSuite {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	o := postgres.NewUserOutbox(db)
	return &UserOutboxSuite{
		outbox: o,
		db:     db,
		tx:     tx,
	}
}
