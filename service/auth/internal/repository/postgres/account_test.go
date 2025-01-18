package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type AccountSuite struct {
	account *postgres.Account
	db      pgxmock.PgxPoolIface
	getter  *mock_uow.MockTxGetter
}

func TestNewAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Account", func(t *testing.T) {
		st := createAccountSuite(t, ctrl)
		assert.NotNil(t, st.account)
	})
}

func TestAccount_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `INSERT INTO accounts \(id, user_id, email, password, created_at, updated_at, created_by, updated_by\) 
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

	t.Run("nil account is prohibited", func(t *testing.T) {
		st := createAccountSuite(t, ctrl)

		err := st.account.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyAccount(), err)
	})

	t.Run("insert duplicate account", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := st.account.Insert(testCtx, account)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			WillReturnError(assert.AnError)

		err := st.account.Insert(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("success insert account", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := st.account.Insert(testCtx, account)

		assert.NoError(t, err)
	})
}

func TestAccount_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `SELECT id, user_id, email, password, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM accounts WHERE email = \$1 LIMIT 1`

	t.Run("get by email returns empty row", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(acc.Email).WillReturnError(pgx.ErrNoRows)

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("get by email returns error", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(acc.Email).WillReturnError(assert.AnError)

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success select by email", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(acc.Email).WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "email", "password", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(acc.ID, acc.UserID, acc.Email, acc.Password, acc.CreatedAt, acc.UpdatedAt, acc.DeletedAt, acc.CreatedBy, acc.UpdatedBy, acc.DeletedBy))

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createTestAccount() *entity.Account {
	return &entity.Account{
		ID:       uuid.Must(uuid.NewV7()),
		UserID:   uuid.Must(uuid.NewV7()),
		Email:    "first@account.com",
		Password: "password",
	}
}

func createAccountSuite(t *testing.T, ctrl *gomock.Controller) *AccountSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	defer pool.Close()
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := uow.NewTxDB(pool, g)
	q := db.New(tx)
	ac := postgres.NewAccount(q)
	return &AccountSuite{
		account: ac,
		db:      pool,
		getter:  g,
	}
}
