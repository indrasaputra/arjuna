package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
)

const (
	queryUpdateRecordStatus = "UPDATE transactions_outbox SET status = ? WHERE id = ?"
)

var (
	testErrInternal = entity.ErrInternal("")
)

type TransactionOutboxSuite struct {
	outbox *postgres.TransactionOutbox
	db     *mock_uow.MockDB
	tx     *mock_uow.MockTx
}

func TestNewTransactionOutbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		assert.NotNil(t, st.outbox)
	})
}

func TestTransactionOutbox_InsertWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := "INSERT INTO " +
		"transactions_outbox (id, status, payload, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil tx is prohibited", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)

		err := st.outbox.InsertWithTx(testCtx, nil, nil)

		assert.Error(t, err)
	})

	t.Run("nil outbox is prohibited", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)

		err := st.outbox.InsertWithTx(testCtx, st.tx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
	})

	t.Run("nil payload is prohibited", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		out := createTestTransactionOutbox()
		out.Payload = nil

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
	})

	t.Run("insert duplicate outbox", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		out := createTestTransactionOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		out := createTestTransactionOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.Error(t, err)
	})

	t.Run("insert returns success", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		out := createTestTransactionOutbox()
		st.tx.EXPECT().
			Exec(testCtx, query, out.ID, out.Status, out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).
			Return(int64(0), nil)

		err := st.outbox.InsertWithTx(testCtx, st.tx, out)

		assert.NoError(t, err)
	})
}

func TestTransactionOutbox_GetAllReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "SELECT id, status, payload FROM transactions_outbox WHERE status = ? ORDER BY created_at ASC LIMIT ? FOR UPDATE"
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, entity.TransactionOutboxStatusReady, limit).
			Return(testErrInternal)

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, entity.TransactionOutboxStatusReady, limit).
			Return(nil)

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestTransactionOutbox_SetProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set processed returns error", func(t *testing.T) {
		out := createTestTransactionOutbox()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusProcessed, out.ID).
			Return(int64(0), testErrInternal)

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set processed success", func(t *testing.T) {
		out := createTestTransaction()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusProcessed, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestTransactionOutbox_SetDelivered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set delivered returns error", func(t *testing.T) {
		out := createTestTransactionOutbox()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusDelivered, out.ID).
			Return(int64(0), testErrInternal)

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set delivered success", func(t *testing.T) {
		out := createTestTransaction()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusDelivered, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestTransactionOutbox_SetFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set failed returns error", func(t *testing.T) {
		out := createTestTransactionOutbox()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusFailed, out.ID).
			Return(int64(0), testErrInternal)

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set failed success", func(t *testing.T) {
		out := createTestTransaction()
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, entity.TransactionOutboxStatusFailed, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestTransactionOutbox_SetRecordStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("set status returns error", func(t *testing.T) {
		out := createTestTransactionOutbox()
		status := entity.TransactionOutboxStatusProcessed
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, status, out.ID).
			Return(int64(0), testErrInternal)

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.Error(t, err)
	})

	t.Run("set status success", func(t *testing.T) {
		out := createTestTransaction()
		status := entity.TransactionOutboxStatusProcessed
		st := createTransactionOutboxSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, queryUpdateRecordStatus, status, out.ID).
			Return(int64(0), nil)

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.NoError(t, err)
	})
}

func createTestTransactionOutbox() *entity.TransactionOutbox {
	transaction := createTestTransaction()
	return &entity.TransactionOutbox{
		ID:      "1",
		Status:  entity.TransactionOutboxStatusReady,
		Payload: transaction,
	}
}

func createTransactionOutboxSuite(ctrl *gomock.Controller) *TransactionOutboxSuite {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	o := postgres.NewTransactionOutbox(db)
	return &TransactionOutboxSuite{
		outbox: o,
		db:     db,
		tx:     tx,
	}
}
