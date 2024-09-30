package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
)

var (
	testCtx = context.Background()
)

type BunDBSuite struct {
	sqldb *sql.DB
	mock  sqlmock.Sqlmock
	bun   *postgres.BunDB
}

func TestNewTxGetter(t *testing.T) {
	t.Run("success create tx getter", func(t *testing.T) {
		g := postgres.NewTxGetter()

		assert.NotNil(t, g)
	})
}

func TestNewTxManager(t *testing.T) {
	t.Run("success create tx manager", func(t *testing.T) {
		tx, err := postgres.NewTxManager(&pgxpool.Pool{})

		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}

func TestNewPgxPool(t *testing.T) {
	t.Run("success create pgx pool", func(t *testing.T) {
		pool, err := postgres.NewPgxPool(postgres.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}

func TestNewDBWithPgx(t *testing.T) {
	t.Run("success create db with pgx", func(t *testing.T) {
		db, err := postgres.NewDBWithPgx(postgres.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, db)
	})
}

func TestNewBunDB(t *testing.T) {
	t.Run("fail create an instance of BunDB", func(t *testing.T) {
		db, err := postgres.NewBunDB(nil)

		assert.Error(t, err)
		assert.Equal(t, postgres.ErrNullDB, err)
		assert.Nil(t, db)
	})

	t.Run("success create an instance of BunDB", func(t *testing.T) {
		db, err := postgres.NewBunDB(&bun.DB{})

		assert.NoError(t, err)
		assert.NotNil(t, db)
	})
}

func TestBunDB_Begin(t *testing.T) {
	t.Run("fail begin the transaction due to unexpected mock call", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()

		tx, err := st.bun.Begin(testCtx)

		assert.Error(t, err)
		assert.Nil(t, tx)
	})

	t.Run("success begin the transaction", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()

		tx, err := st.bun.Begin(testCtx)

		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}

func TestBunDB_Query(t *testing.T) {
	query := `SELECT id FROM tables`

	t.Run("query returns error", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		var dest interface{}
		errReturn := errors.New("error")

		st.mock.ExpectQuery(query).WillReturnError(errReturn)

		err := st.bun.Query(testCtx, &dest, query)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("query success", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		var id int

		st.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := st.bun.Query(testCtx, &id, query)

		assert.NoError(t, err)
	})
}

func TestBunTx_Begin(t *testing.T) {
	t.Run("tx begin return itself", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		res, _ := st.bun.Begin(testCtx)

		tx, err := res.Begin(testCtx)

		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, res, tx)
	})
}

func TestBunTx_Commit(t *testing.T) {
	t.Run("commit returns error", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		errReturn := errors.New("error")
		st.mock.ExpectCommit().WillReturnError(errReturn)

		err := tx.Commit(testCtx)

		assert.Error(t, err)
	})

	t.Run("commit success", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		st.mock.ExpectCommit()

		err := tx.Commit(testCtx)

		assert.NoError(t, err)
	})
}

func TestBunTx_Rollback(t *testing.T) {
	t.Run("rollback returns error", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		errReturn := errors.New("error")
		st.mock.ExpectRollback().WillReturnError(errReturn)

		err := tx.Rollback(testCtx)

		assert.Error(t, err)
	})

	t.Run("rollback success", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		st.mock.ExpectRollback()

		err := tx.Rollback(testCtx)

		assert.NoError(t, err)
	})
}

func TestBunTx_Query(t *testing.T) {
	query := `SELECT id FROM tables`

	t.Run("query returns error", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		errReturn := errors.New("error")
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		st.mock.ExpectQuery(query).WillReturnError(errReturn)
		var dest interface{}

		err := tx.Query(testCtx, &dest, query)

		assert.Error(t, err)
		assert.Equal(t, errReturn, err)
	})

	t.Run("query success", func(t *testing.T) {
		st := createBunDBSuite(t)
		defer func() {
			_ = st.sqldb.Close()
		}()
		st.mock.ExpectBegin()
		tx, _ := st.bun.Begin(testCtx)
		st.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		var id int

		err := tx.Query(testCtx, &id, query)

		assert.NoError(t, err)
	})
}

func createBunDBSuite(t *testing.T) *BunDBSuite {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	bdb := bun.NewDB(db, pgdialect.New())
	b, _ := postgres.NewBunDB(bdb)

	return &BunDBSuite{
		sqldb: db,
		mock:  mock,
		bun:   b,
	}
}
