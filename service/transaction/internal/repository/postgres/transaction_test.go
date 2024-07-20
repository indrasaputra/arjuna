package postgres_test

import (
	"context"
	"testing"

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
	tx  *mock_uow.MockTx
}

func TestNewTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransaction_InsertWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "INSERT INTO " +
		"transactions (id, sender_id, receiver_id, amount, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil tx is prohibited", func(t *testing.T) {
		st := createTransactionSuite(ctrl)

		err := st.trx.InsertWithTx(testCtx, nil, nil)

		assert.Error(t, err)
	})

	t.Run("nil transactions is prohibited", func(t *testing.T) {
		st := createTransactionSuite(ctrl)

		err := st.trx.InsertWithTx(testCtx, st.tx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
	})

	t.Run("insert duplicate transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.trx.InsertWithTx(testCtx, st.tx, trx)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.trx.InsertWithTx(testCtx, st.tx, trx)

		assert.Error(t, err)
	})

	t.Run("success insert transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			Return(int64(1), nil)

		err := st.trx.InsertWithTx(testCtx, st.tx, trx)

		assert.NoError(t, err)
	})
}

func createTestTransaction() *entity.Transaction {
	a, _ := decimal.NewFromString("10.23")
	return &entity.Transaction{
		ID:         "123",
		SenderID:   "1",
		ReceiverID: "2",
		Amount:     a,
	}
}

func createTransactionSuite(ctrl *gomock.Controller) *TransactionSuite {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	t := postgres.NewTransaction(db)
	return &TransactionSuite{
		trx: t,
		db:  db,
		tx:  tx,
	}
}
