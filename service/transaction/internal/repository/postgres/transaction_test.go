package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type TransactionSuite struct {
	trx    *postgres.Transaction
	db     pgxmock.PgxPoolIface
	getter *mock_uow.MockTxGetter
}

func TestNewTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionSuite(t, ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransaction_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `INSERT INTO transactions \(id, sender_id, receiver_id, amount, created_at, updated_at, created_by, updated_by\)
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

	t.Run("nil transactions is prohibited", func(t *testing.T) {
		st := createTransactionSuite(t, ctrl)

		err := st.trx.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
	})

	t.Run("insert duplicate transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := st.trx.Insert(testCtx, trx)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			WillReturnError(assert.AnError)

		err := st.trx.Insert(testCtx, trx)

		assert.Error(t, err)
	})

	t.Run("success insert transactions", func(t *testing.T) {
		trx := createTestTransaction()
		st := createTransactionSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(trx.ID, trx.SenderID, trx.ReceiverID, trx.Amount, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

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
		st := createTransactionSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WillReturnError(assert.AnError)

		err := st.trx.DeleteAll(testCtx)

		assert.Error(t, err)
	})

	t.Run("delete all returns success", func(t *testing.T) {
		st := createTransactionSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WillReturnResult(pgxmock.NewResult("DELETE", 1))

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

func createTransactionSuite(t *testing.T, ctrl *gomock.Controller) *TransactionSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	defer pool.Close()
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := sdkpostgres.NewTxDB(pool, g)
	q := db.New(tx)
	trx := postgres.NewTransaction(q)
	return &TransactionSuite{
		trx:    trx,
		db:     pool,
		getter: g,
	}
}
