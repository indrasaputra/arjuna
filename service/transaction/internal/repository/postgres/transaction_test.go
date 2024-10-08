package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type TransactionSuite struct {
	trx *postgres.Transaction
	db  *mock_uow.MockDB
}

func TestNewTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransaction_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "INSERT INTO " +
		"transactions (id, sender_id, receiver_id, amount, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil transactions is prohibited", func(t *testing.T) {
		st := createTransactionSuite(ctrl)

		err := st.trx.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
	})

	t.Run("insert duplicate transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.trx.Insert(testCtx, trx)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.trx.Insert(testCtx, trx)

		assert.Error(t, err)
	})

	t.Run("success insert transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(1), nil)

		err := st.trx.Insert(testCtx, trx)

		assert.NoError(t, err)
	})
}

func TestTransaction_DeleteAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `DELETE FROM transactions`

	t.Run("delete all returns error", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query).
			Return(int64(0), entity.ErrInternal("error"))

		err := st.trx.DeleteAll(testCtx)

		assert.Error(t, err)
	})

	t.Run("delete all returns success", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query).
			Return(int64(1), nil)

		err := st.trx.DeleteAll(testCtx)

		assert.NoError(t, err)
	})
}

func createTestTransaction() *entity.Transaction {
	a, _ := decimal.NewFromString("10.23")
	return &entity.Transaction{
		ID:         uuid.Must(uuid.NewV7()),
		SenderID:   uuid.Must(uuid.NewV7()),
		ReceiverID: uuid.Must(uuid.NewV7()),
		Amount:     a,
	}
}

func createTransactionSuite(ctrl *gomock.Controller) *TransactionSuite {
	db := mock_uow.NewMockDB(ctrl)
	t := postgres.NewTransaction(db)
	return &TransactionSuite{
		trx: t,
		db:  db,
	}
}
