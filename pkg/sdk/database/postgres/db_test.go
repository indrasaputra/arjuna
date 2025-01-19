package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
)

var (
	testCtx = context.Background()
)

type TxDBSuite struct {
	db     pgxmock.PgxPoolIface
	getter *mock_uow.MockTxGetter
	tx     *postgres.TxDB
}

func TestIsUniqueViolationError(t *testing.T) {
	t.Run("error is not unique violation", func(t *testing.T) {
		res := postgres.IsUniqueViolationError(assert.AnError)
		assert.False(t, res)
	})

	t.Run("nil error is not unique violation", func(t *testing.T) {
		res := postgres.IsUniqueViolationError(nil)
		assert.False(t, res)
	})

	t.Run("error is unique violation", func(t *testing.T) {
		res := postgres.IsUniqueViolationError(&pgconn.PgError{Code: "23505"})
		assert.True(t, res)
	})
}

func TestNewPgxPool(t *testing.T) {
	t.Run("success create pgx pool", func(t *testing.T) {
		pool, err := postgres.NewPgxPool(postgres.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}

func TestNewTxDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success create txdb", func(t *testing.T) {
		st := createTxDBSuite(t, ctrl)

		assert.NotNil(t, st.tx)
	})
}

func TestTxDB_Exec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success execute exec", func(t *testing.T) {
		st := createTxDBSuite(t, ctrl)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec("exec").WithArgs("arg1", "arg2").WillReturnResult(pgxmock.NewResult("exec", 1))

		res, err := st.tx.Exec(testCtx, "exec", "arg1", "arg2")

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestTxDB_Query(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success execute query", func(t *testing.T) {
		st := createTxDBSuite(t, ctrl)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery("query").
			WithArgs("arg1", "arg2").
			WillReturnRows(pgxmock.
				NewRows([]string{"id"}).
				AddRow("id"))

		res, err := st.tx.Query(testCtx, "query", "arg1", "arg2")

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestTxDB_QueryRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success execute query row", func(t *testing.T) {
		st := createTxDBSuite(t, ctrl)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery("query").
			WithArgs("arg1", "arg2").
			WillReturnRows(pgxmock.
				NewRows([]string{"id"}).
				AddRow("id"))

		res := st.tx.QueryRow(testCtx, "query", "arg1", "arg2")

		assert.NotNil(t, res)
	})
}

func createTxDBSuite(t *testing.T, ctrl *gomock.Controller) *TxDBSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := postgres.NewTxDB(pool, g)
	return &TxDBSuite{
		db:     pool,
		getter: g,
		tx:     tx,
	}
}
